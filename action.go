package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/monochromegane/go-gitignore"
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

// Walker is a struct for walking dir
type Walker struct {
	found     []string
	gitignore gitignore.IgnoreMatcher
	ignore    *regexp.Regexp
	scanned   int
	skipped   int
	words     [][]byte
}

// NewWalker is a constructor of Walker
func NewWalker() (w *Walker, err error) {
	matcher, err := Config.GitIgnore()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	w = &Walker{
		gitignore: matcher,
		words:     make([][]byte, len(Config.Words)),
	}
	for i, word := range Config.Words {
		w.words[i] = []byte(word)
	}
	ignores := strings.Join(Config.Ignores, "|")
	if w.ignore, err = regexp.Compile(ignores); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return
}

func (w *Walker) walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	w.scanned++
	if w.ignore.MatchString(path) {
		w.skipped++
		return nil
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return xerrors.Errorf("reading from %s: %w", path, err)
	}
	contained := false
	for _, word := range w.words {
		if bytes.Contains(content, word) {
			contained = true
			break
		}
	}
	if contained {
		w.found = append(w.found, path)
	}
	w.scanned++
	return nil
}

// ResultHeader returns header string for the result
func (w *Walker) ResultHeader() string {
	return fmt.Sprintf(
		"dotfiles scan for prohibited words; scanned: %d files(s), skipped: %d files(s), found: %d files(s)",
		w.scanned, w.skipped, len(w.found))
}
