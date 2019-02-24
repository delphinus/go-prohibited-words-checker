package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
	"golang.org/x/xerrors"
)

// Walker is a struct for walking dir
type Walker struct {
	found     []string
	gitignore *ignore.GitIgnore
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
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if info.IsDir() {
		if err := w.walkDir(path); err != nil {
			if err == filepath.SkipDir {
				return err
			}
			return xerrors.Errorf(": %w", err)
		}
		return nil
	}
	if err := w.walkFile(path); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (w *Walker) walkDir(path string) error {
	rel, err := filepath.Rel(Config.Dir, path)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	// sabhiram/go-gitignore cannot ignore directories without trailing
	// slash for the .gitignore entries with tariling slashes. such as...
	//
	//   in .gitignore:
	//     /node_modules/
	//   the directory:
	//     node_modules
	if !strings.HasSuffix(rel, "/") {
		rel += "/"
	}
	if w.gitignore.MatchesPath(rel) {
		return filepath.SkipDir
	}
	return nil
}

func (w *Walker) walkFile(path string) error {
	rel, err := filepath.Rel(Config.Dir, path)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if w.gitignore.MatchesPath(rel) || w.ignore.MatchString(rel) {
		w.scanned++
		w.skipped++
		return nil
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return xerrors.Errorf("reading from %s: %w", rel, err)
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

// Data returns a map to use in templates
func (w *Walker) Data() map[string]interface{} {
	return map[string]interface{}{}
}
