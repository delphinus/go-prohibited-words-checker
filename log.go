package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

// Log takes the same args as Sprintf and print logs if on the `verbose`
// flag.
func Log(cli *cli.Context, format string, a ...interface{}) {
	if cli.Bool("verbose") {
		fmt.Fprintf(cli.App.Writer, format, a...)
	}
}
