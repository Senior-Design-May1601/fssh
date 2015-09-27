package main

import (
    "errors"
    "flag"
    "io/ioutil"
    "net"
    "strconv"

    "golang.org/x/crypto/ssh"
)

const (
    DEFAULT_KEY = "../keys/dummy_id_rsa"
    DEFAULT_PORT = 2222
)

func readSecretKey(path string) (ssh.Signer) {
    raw, err := ioutil.ReadFile(path)
    if err != nil {
        panic(err)
    }

    sk, err := ssh.ParsePrivateKey(raw)
    if err != nil {
        panic(err)
    }

    return sk
}

func keyHandler(c ssh.ConnMetadata, k ssh.PublicKey) (*ssh.Permissions, error) {
    return nil, errors.New("")
}

func passwdHandler(c ssh.ConnMetadata, p []byte) (*ssh.Permissions, error) {
    return nil, errors.New("")
}

func main() {
    p := flag.Int("port", DEFAULT_PORT, "SSH server port")
    key := flag.String("key", DEFAULT_KEY, "path to SSH private key")
    flag.Parse()
    port := ":" + strconv.Itoa(*p)

    config := ssh.ServerConfig{
        PublicKeyCallback: keyHandler,
        PasswordCallback: passwdHandler,
    }
    config.AddHostKey(readSecretKey(*key))

    s, err := net.Listen("tcp", port)
    if err != nil {
        panic(err)
    }
    defer s.Close()

    for {
        c, err := s.Accept()
        if err != nil {
            panic(err)
        }

        ssh.NewServerConn(c, &config)
    }
}
