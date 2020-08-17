package internal

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func Parse(file, arg string) ([]target, []string) {
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

	err = viper.UnmarshalKey(arg, &r)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range r.Target {
		var item target
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
