package internal

import (
	"github.com/spf13/viper"
	"log"
	"os"
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
	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(f)
	if err != nil {
		log.Fatalln(err)
	}
}
