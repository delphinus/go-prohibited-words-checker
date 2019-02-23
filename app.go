package main

import "gopkg.in/urfave/cli.v2"

// NewApp will return the App
func NewApp() *cli.App {
	return &cli.App{
		Usage:   "Check for prohibited words",
		Version: Version,
		Authors: []*cli.Author{
			{Name: "JINNOUCHI Yasushi", Email: "delphinus@remora.cx"},
		},
		Before: LoadConfig,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Show logs verbosely",
			},
		},
	}
}
