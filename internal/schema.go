package internal

import "golang.org/x/crypto/ssh"

type target struct {
	User          string
	Address       string
	Port          int
	PrivateKey    string `mapstructure:"private_key"`
	KnownHosts    string `mapstructure:"known_hosts"`
	IgnoreHostKey bool   `mapstructure:"ignore_host_key"`
}

type run struct {
	Target []string
	Script []string
}

type config struct {
	Global target
}

type client struct {
	*ssh.Client
}
