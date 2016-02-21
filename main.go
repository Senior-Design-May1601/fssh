package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/Senior-Design-May1601/projectmain/logger"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	Address string
	Port    int
	Key     string
}

func readSecretKey(path string) ssh.Signer {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		mylogger.Fatal(err)
	}

	sk, err := ssh.ParsePrivateKey(raw)
	if err != nil {
		mylogger.Fatal(err)
	}

	return sk
}

func keyHandler(c ssh.ConnMetadata, k ssh.PublicKey) (*ssh.Permissions, error) {
	mylogger.Println("key login attempt")
	return nil, errors.New("")
}

func passwdHandler(c ssh.ConnMetadata, p []byte) (*ssh.Permissions, error) {
	mylogger.Println("password login attempt")
	return nil, errors.New("")
}

var mylogger *log.Logger

func main() {
	mylogger = logger.NewLogger("", 0)

	configPath := flag.String("config", "", "fssh config file")
	flag.Parse()

	var config Config
	if _, err := toml.DecodeFile(*configPath, &config); err != nil {
		mylogger.Fatal(err)
	}

	sshConfig := ssh.ServerConfig{
		PublicKeyCallback: keyHandler,
		PasswordCallback:  passwdHandler,
	}
	sshConfig.AddHostKey(readSecretKey(config.Key))

	s, err := net.Listen("tcp", config.Address+":"+strconv.Itoa(config.Port))
	if err != nil {
		mylogger.Fatal(err)
	}
	defer s.Close()

	for {
		c, err := s.Accept()
		if err != nil {
			mylogger.Fatal(err)
		}

		ssh.NewServerConn(c, &sshConfig)
	}
}
