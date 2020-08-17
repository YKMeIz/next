package main

import (
	"errors"
	"github.com/YKMeIz/next/internal"
	"github.com/spf13/cobra"
	"log"
	"os"
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
			Use:   "hugo",
			Short: "Hugo is a very fast static site generator",
			Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
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
			return errors.New("one or more arguments are required")
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
