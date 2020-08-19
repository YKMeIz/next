package internal

import (
	"os/user"
	"reflect"
	"testing"
)

func TestParseVariableOverride(t *testing.T) {
	targets, cmds := Parse("../examples/variable_override.yml", "run_uptime")

	if cmds[0] != "uptime" {
		t.Error("Parse() gets wrong script:\n expect:\n", "uptime", "\n get:\n", cmds[0])
	}

	if i := (target{
		User:          "ykmeiz",
		Address:       "10.0.0.155",
		Port:          22,
		PrivateKey:    "/Users/ykmeiz/.ssh/id_rsa",
		KnownHosts:    "/Users/ykmeiz/.ssh/known_hosts",
		IgnoreHostKey: true,
	}); targets[0] != i {
		t.Error("Parse() gets wrong target:\n expect:\n", i, "\n get:\n", targets[0])
	}

	if i := (target{
		User:          "ykmeiz",
		Address:       "10.0.0.156",
		Port:          33,
		PrivateKey:    "/home/ykmeiz/.ssh/id_rsa",
		KnownHosts:    "/home/ykmeiz/.ssh/known_hosts",
		IgnoreHostKey: false,
	}); targets[1] != i {
		t.Error("Parse() gets wrong target:\n expect:\n", i, "\n get:\n", targets[1])
	}
}

func TestParseDefault(t *testing.T) {
	targets, cmds := Parse("../examples/default.yml", "run_uname")

	if c := []string{"echo \"This is $(hostname):\"", "uname -a"}; !reflect.DeepEqual(cmds, c) {
		t.Error("Parse() gets wrong script:\n expect:\n", c, "\n get:\n", cmds)
	}

	u, err := user.Current()
	if err != nil {
		t.Error(err)
	}

	if i := (target{
		User:          u.Username,
		Address:       "10.0.0.155",
		Port:          22,
		PrivateKey:    u.HomeDir + "/.ssh/id_rsa",
		KnownHosts:    u.HomeDir + "/.ssh/known_hosts",
		IgnoreHostKey: false,
	}); targets[0] != i {
		t.Error("Parse() gets wrong target:\n expect:\n", i, "\n get:\n", targets[0])
	}

}
