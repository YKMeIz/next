package internal

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"strconv"
)

// Connect makes SSH connection and return client.
func Connect(t target) *client {
	// Read private key.
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
		User: t.User,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Host key needs to be checked.
	if !t.IgnoreHostKey {
		pubKey, err := ioutil.ReadFile(t.KnownHosts)
		if err != nil {
			log.Fatalf("unable to read public key: %v", err)
		}

		var hostKey ssh.PublicKey

		// Loop all hosts in file and find the key matches given address.
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

// RemoteExec creates SSH session and executes given command.
func (c *client) RemoteExec(cmd string) string {
	session, err := c.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(cmd); err != nil {
		log.Println(err.Error())
		return stderr.String()
	}

	return stdout.String()
}
