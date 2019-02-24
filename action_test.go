package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestNewWalk(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		re []string
		ok bool
	}{
		{re: []string{"invalid[re"}, ok: false},
		{re: []string{"valid[re]"}, ok: true},
	} {
		(func() {
			defer prepareValidConfig(t)()
			Config.Ignores = c.re
			if _, err := NewWalker(); c.ok {
				a.NoError(err)
			} else {
				a.Error(err)
				t.Logf("found err: %s", err)
			}
		})()
	}
}

func TestIgnore(t *testing.T) {
	a := assert.New(t)
	prepareValidConfig(t)
	gi, err := Config.GitIgnore()
	a.NoError(err)
	_ = createFile(t, ".vim/vimrc", "")
	a.False(gi.MatchesPath("hoge"))
	a.True(gi.MatchesPath(".vim"))
	a.True(gi.MatchesPath(".vim/"))
	a.True(gi.MatchesPath("hoge/.vim/"))
	a.True(gi.MatchesPath(".vim/vimrc"))
	a.False(gi.MatchesPath("node_modules"))
	a.True(gi.MatchesPath("node_modules/"))
	a.False(gi.MatchesPath("hoge/node_modules/"))
}

func TestWalkWithError(t *testing.T) {
	a := assert.New(t)
	defer prepareValidConfig(t)()
	w, err := NewWalker()
	a.NoError(err)
	argErr := xerrors.New("hogehogeo")
	err = w.walk("", nil, argErr)
	a.Error(err)
	a.True(xerrors.Is(err, argErr))
}

func TestWalkDir(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		path    string
		skipDir bool
	}{
		{path: ".vim", skipDir: true},
		{path: ".vim/", skipDir: true},
		{path: "node_modules", skipDir: true},
		{path: "node_modules/", skipDir: true},
		{path: "some-dir", skipDir: false},
	} {
		(func() {
			defer prepareValidConfig(t)()
			w, err := NewWalker()
			a.NoError(err, c.path)
			path := filepath.Join(Config.Dir, c.path)
			// directory os.FileInfo is needed with tests
			info, err := os.Stat(Config.Dir)
			a.NoError(err, c.path)
			err = w.walk(path, info, nil)
			a.Equal(0, w.scanned, c.path)
			a.Equal(0, w.skipped, c.path)
			a.Len(w.found, 0, c.path)
			if c.skipDir {
				a.Error(err, c.path)
				a.Equal(filepath.SkipDir, err, c.path)
			} else {
				a.NoError(err, c.path)
			}
		})()
	}
}

func TestWalkFileSkipped(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		path    string
		scanned int
		skipped int
	}{
		// ignored by .gitignore
		{path: ".vim/vimrc", scanned: 1, skipped: 1},
		{path: "node_modules/hoge", scanned: 1, skipped: 1},
		{path: "node_modules/hoge/fuga", scanned: 1, skipped: 1},
		// ignored by regexp
		{path: "want_to_ignore/hogehoge", scanned: 1, skipped: 1},
		// not ignored
		{path: ".git/config", scanned: 1, skipped: 0},
		{path: "hoge/node_modules/fuga", scanned: 1, skipped: 0},
	} {
		(func() {
			defer prepareValidConfig(t)()
			w, err := NewWalker()
			a.NoError(err)
			path := createFile(t, c.path, "")
			info, err := os.Stat(path)
			a.NoError(err)
			err = w.walk(path, info, nil)
			a.Equal(c.scanned, w.scanned, c.path+" for scanned")
			a.Equal(c.skipped, w.skipped, c.path+" for skipped")
			a.NoError(err, c.path+" for error")
		})()
	}
}

func TestWalkFile(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		content string
		found   bool
	}{
		{content: "HOGE", found: true},
		{content: "FUGA", found: true},
		{content: "hogeHOGEo", found: true},
		{content: "hoge", found: false},
	} {
		(func() {
			defer prepareValidConfig(t)()
			w, err := NewWalker()
			a.NoError(err, c.content)
			file := "hoge.txt"
			path := createFile(t, file, c.content)
			info, err := os.Stat(path)
			a.NoError(err, c.content)
			err = w.walk(path, info, nil)
			a.NoError(err, c.content)
			a.Equal(1, w.scanned, c.content)
			a.Equal(0, w.skipped, c.content)
			if c.found {
				a.Len(w.found, 1, c.content)
				a.Equal(path, w.found[0], c.content)
			} else {
				a.Len(w.found, 0, c.content)
			}
		})()
	}
}

func createFile(t *testing.T, file string, content string) string {
	a := assert.New(t)
	path := filepath.Join(Config.Dir, file)
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		a.NoError(os.MkdirAll(dir, 0700))
	}
	a.NoError(ioutil.WriteFile(path, []byte(content), 0600))
	return path
}
