package internal

import (
	"log"
	"os/user"
	"reflect"
)

func targetCompletion(global, local target) target {
	t := merge(global, local)

	user, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	if t.User == "" {
		t.User = user.Username
	}

	if t.Port == 0 {
		t.Port = 22
	}

	if t.PrivateKey == "" {
		t.PrivateKey = user.HomeDir + "/.ssh/id_rsa"
	}

	if t.KnownHosts == "" && !t.IgnoreHostKey {
		t.KnownHosts = user.HomeDir + "/.ssh/known_hosts"
	}

	return t
}

func merge(global, local target) target {
	l := reflect.ValueOf(&local).Elem()
	g := reflect.ValueOf(&global).Elem()

	for i := 0; i < l.NumField(); i++ {
		if !l.Field(i).IsZero() {
			g.Field(i).Set(l.Field(i))
		}
	}

	return g.Interface().(target)
}
