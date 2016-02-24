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
	"github.com/Senior-Design-May1601/Splunk/alert"
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

func baseAlertMap(metadata ssh.ConnMetadata) map[string]string {
    meta := make(map[string]string)
    meta["service"] = "ssh"
    meta["user"] = metadata.User()
    meta["remote"] = metadata.RemoteAddr().String()
    meta["local"] = metadata.LocalAddr().String()

    return meta
}

func keyAlert(metadata ssh.ConnMetadata, key ssh.PublicKey) string {
    meta := baseAlertMap(metadata)
    meta["authtype"] = "publickey"
    meta["key"] = string(ssh.MarshalAuthorizedKey(key))

    return alert.NewSplunkAlertMessage(meta)
}

func passwdAlert(metadata ssh.ConnMetadata, passwd []byte) string {
    meta := baseAlertMap(metadata)
    meta["authtype"] = "password"
    meta["password"] = string(passwd)

    return alert.NewSplunkAlertMessage(meta)
}

func keyHandler(meta ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	mylogger.Println(keyAlert(meta, key))
	return nil, errors.New("")
}

func passwdHandler(meta ssh.ConnMetadata, p []byte) (*ssh.Permissions, error) {
	mylogger.Println(passwdAlert(meta, p))
	return nil, errors.New("")
}

var mylogger *log.Logger
var config Config

func main() {
	mylogger = logger.NewLogger("", 0)

	configPath := flag.String("config", "", "fssh config file")
	flag.Parse()

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
