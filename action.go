package main

import (
	"fmt"
	"path/filepath"

	"golang.org/x/xerrors"
	"gopkg.in/urfave/cli.v2"
)

// Action is the main logic
func Action(cli *cli.Context) error {
	w, err := NewWalker()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := filepath.Walk(Config.Dir, w.walk); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	Log(cli, w.ResultHeader())
	return nil
}

// ResultHeader returns header string for the result
func (w *Walker) ResultHeader() string {
	return fmt.Sprintf(
		"dotfiles scan for prohibited words; scanned: %d files(s), "+
			"skipped: %d files(s), found: %d files(s)",
		w.scanned, w.skipped, len(w.found))
}
