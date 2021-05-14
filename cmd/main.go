package main

import (
	"errors"
	"github.com/YKMeIz/next/internal"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"sync"
)

type cli struct {
	*cobra.Command

	conf string
	wg   sync.WaitGroup
}

func newCmd() *cli {
	c := &cli{
		Command: &cobra.Command{
			Use:   "next",
			Short: "IT automation tool enabling remote commands execution on multiple machines.",
		},
		conf: "next.yml",
	}

	c.RunE = c.run()
	c.Flags().StringVarP(&c.conf, "config", "c", c.conf, "specify a configuration file")

	return c
}

func (c *cli) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			availableCmds := internal.ParseAvailableCommands(c.conf)
			for _, v := range availableCmds {
				c.AddCommand(&cobra.Command{Use: v})
			}

			if len(availableCmds) == 0 {
				return errors.New("configuration file is not set, or no argument found in configuration file\n")
			}

			return errors.New("one or more arguments are required.\n\nAvailable arguments are: " + strings.Join(availableCmds, ", ") + "\n")
		}

		for _, v := range args {
			c.wg.Add(1)
			go c.exec(v)
		}

		c.wg.Wait()

		return nil
	}
}

func (c *cli) exec(arg string) {
	t, s := internal.Parse(c.conf, arg)
	for _, v := range t {
		c := internal.Connect(v)
		for _, cmd := range s {
			log.Print(v.Address, ": ", c.RemoteExec(cmd))
		}
	}
	c.wg.Done()
}

func main() {
	if err := newCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
