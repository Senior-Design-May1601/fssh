# fssh
A fake SSH server.

### Description

`fssh` opens a socket and listens for incoming SSH connections. `fssh` steps
through the initial SSH handshake and accepts a password or public key,
always rejecting the login attempt.

### Usage
```
Usage of /tmp/go-build564020242/command-line-arguments/_obj/exe/fssh:
  -key="../keys/dummy_id_rsa": path to SSH private key
  -port=2222: SSH server port
```
