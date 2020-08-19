package internal

import (
	"log"
	"os/user"
	"reflect"
)

// targetCompletion generates target details.
// Global variable will be overridden if local variable is defined.
// username, private key, known hosts, and SSH port will be set as system default if they are not defined by user.
func targetCompletion(global, local target) target {
	t := merge(global, local)

	user, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	// Set current username if username is undefined.
	if t.User == "" {
		t.User = user.Username
	}

	// Set default SSH port if port is undefined.
	if t.Port == 0 {
		t.Port = 22
	}

	// Set private key file location if the key is undefined.
	if t.PrivateKey == "" {
		t.PrivateKey = user.HomeDir + "/.ssh/id_rsa"
	}

	// Set known hosts file location if the key is undefined.
	if t.KnownHosts == "" && !t.IgnoreHostKey {
		t.KnownHosts = user.HomeDir + "/.ssh/known_hosts"
	}

	return t
}

// merge overrides global variable.
func merge(global, local target) target {
	l := reflect.ValueOf(&local).Elem()
	g := reflect.ValueOf(&global).Elem()

	// Loop each struct filed and override value if local variable is non-zero.
	for i := 0; i < l.NumField(); i++ {
		if !l.Field(i).IsZero() {
			g.Field(i).Set(l.Field(i))
		}
	}

	return g.Interface().(target)
}
