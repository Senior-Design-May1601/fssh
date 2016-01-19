package main

import (
    "errors"
    "io/ioutil"
    "log"
    "net"

    "golang.org/x/crypto/ssh"
    "github.com/Senior-Design-May1601/projectmain/plugin"
)

const (
    KEY_PATH = "/path/to/private/key"
    SSH_PORT = "8022"
)

func readSecretKey(path string) (ssh.Signer) {
    raw, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatal(err)
    }

    sk, err := ssh.ParsePrivateKey(raw)
    if err != nil {
        log.Fatal(err)
    }

    return sk
}

func keyHandler(c ssh.ConnMetadata, k ssh.PublicKey) (*ssh.Permissions, error) {
    log.Println("SSH public key connection attempt.")
    return nil, errors.New("")
}

func passwdHandler(c ssh.ConnMetadata, p []byte) (*ssh.Permissions, error) {
    return nil, errors.New("")
}

type SSHPlugin struct {
    port int
}

func (x *SSHPlugin) Start(args *plugin.Args, reply *plugin.Reply) error {
    log.Println("ssh start called")
    sshConfig = ssh.ServerConfig{
        PublicKeyCallback: keyHandler,
        PasswordCallback: passwdHandler,
    }
    sshConfig.AddHostKey(readSecretKey(KEY_PATH))

    listener, err := net.Listen("tcp", "localhost:"+SSH_PORT)
    if err != nil {
        return errors.New("Failed to start SSH server.")
    }
    log.Println("SSH server listening on", SSH_PORT)

    go func() {
        for {
            conn, err := listener.Accept()
            if err != nil {
                log.Fatal(err)
            }
            ssh.NewServerConn(conn, &sshConfig)
        }
    }()

    return nil
}

func (x *SSHPlugin) Stop(args *plugin.Args, reply *plugin.Reply) error {
    log.Println("ssh stop called")
    // TODO: will this actually stop the server?!?
    listener.Close()
    return nil
}

func (x *SSHPlugin) Restart(args *plugin.Args, reply *plugin.Reply) error {
    log.Println("ssh restart called")
    return nil
}

func (x *SSHPlugin) Port() int {
    return 10000
}

var sshConfig ssh.ServerConfig
var listener net.Listener

func main () {
    sshPlugin := &SSHPlugin{port: 10000}
    server, err := plugin.NewPlugin(sshPlugin)
    if err != nil {
        log.Fatal(err)
    }
    server.Serve()
    log.Println("ssh server dying")
}
