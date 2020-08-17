package internal

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"strconv"
)

func Connect(t target) *client {
	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile(t.PrivateKey)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		//User: user.Username,
		User: t.User,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if !t.IgnoreHostKey {
		pubKey, err := ioutil.ReadFile(t.KnownHosts)
		if err != nil {
			log.Fatalf("unable to read public key: %v", err)
		}

		var hostKey ssh.PublicKey

		for {
			_, hosts, key, _, rest, err := ssh.ParseKnownHosts(pubKey)
			if err != nil {
				log.Fatalf("unable to parse public key: %v", err)
			}
			if isIn(t.Address, hosts) {
				hostKey = key
				break
			}
			pubKey = rest
			if len(rest) == 0 {
				log.Fatalf("unable to find host key")
			}
		}

		config.HostKeyCallback = ssh.FixedHostKey(hostKey)
	}

	c, err := ssh.Dial("tcp", t.Address+":"+strconv.Itoa(t.Port), config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	return &client{c}
}

func isIn(s string, search []string) bool {
	for _, v := range search {
		if s == v {
			return true
		}
	}
	return false
}

func (c *client) RemoteExec(cmd string) string {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := c.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(cmd); err != nil {
		log.Println(err.Error())
		return stderr.String()
	}

	return stdout.String()
}
