package internal

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

// Parse unmarshals rule from file, and return target details and commands need to be executed remotely.
func Parse(file, rule string) ([]target, []string) {
	var (
		c config
		r run
		t []target
	)

	readConfig(file)

	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalln(err)
	}

	err = viper.UnmarshalKey(rule, &r)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range r.Target {
		var item target

		// Check if target variable is defined.
		if viper.IsSet(v) {
			if err := viper.UnmarshalKey(v, &item); err != nil {
				log.Fatalln(err)
			}
		} else {
			item.Address = v
		}

		t = append(t, targetCompletion(c.Global, item))
	}

	return t, r.Script
}

func readConfig(file string) {
	f, _ := os.Open(file)

	defer f.Close()

	viper.SetConfigType("yaml")
	_ = viper.ReadConfig(f)
}

// ParseAvailableCommands finds all defined commands as next arguments.
func ParseAvailableCommands(file string) []string {
	readConfig(file)
	var cmds []string
	for _, v := range viper.AllKeys() {
		if strings.Contains(v, "script") {
			cmds = append(cmds, strings.TrimSuffix(v, ".script"))
		}
	}
	return cmds
}
