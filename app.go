package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

// NewApp will return the App
func NewApp() *cli.App {
	return &cli.App{
		Usage:   "Check for prohibited words",
		Version: Version,
		Authors: []*cli.Author{
			{Name: "JINNOUCHI Yasushi", Email: "delphinus@remora.cx"},
		},
		Before: handleExit(LoadConfig),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Show logs verbosely",
			},
			&cli.BoolFlag{
				Name:    "mail",
				Aliases: []string{"m"},
				Usage:   "Send mail for results",
			},
		},
		Action: handleExit(Action),
	}
}

func handleExit(handler func(*cli.Context) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		if err := handler(c); err != nil {
			return cli.Exit(fmt.Sprintf("%+v", err), 1)
		}
		return nil
	}
}
