package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
	"gopkg.in/urfave/cli.v2"
)

func TestGitIgnore(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		hasGitignore bool
		ok           bool
	}{
		{hasGitignore: true, ok: true},
		{hasGitignore: false, ok: false},
	} {
		(func() {
			tmpDir, done := prepare(t)
			defer done()
			if c.hasGitignore {
				a.NoError(
					ioutil.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(""), 0600),
				)
			}
			config := &Configs{Dir: tmpDir}
			if _, err := config.GitIgnore(); c.ok {
				a.NoError(err)
			} else {
				a.Error(err)
			}
		})()
	}
}

var userCurrentErr error

func TestConfigFilename(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		err error
	}{
		{err: nil},
		{err: xerrors.New("some error")},
	} {
		userCurrentErr = c.err
		(func() {
			_, done := prepare(t)
			defer done()
			_, err := configFilename()
			if c.err == nil {
				a.NoError(err)
			} else {
				a.True(xerrors.Is(err, c.err))
			}
		})()
		userCurrentErr = nil
	}
}

func TestNoConfig(t *testing.T) {
	a := assert.New(t)
	_, done := prepare(t)
	defer done()
	err := LoadConfig(&cli.Context{})
	t.Logf("found err: %s", err)
	a.Error(err)
}

func TestInValidConfig(t *testing.T) {
	a := assert.New(t)
	tmpDir, done := prepare(t)
	defer done()
	prepareConfig(t, tmpDir, []byte(`
invalid hogehoge
`))
	err := LoadConfig(&cli.Context{})
	t.Logf("found err: %s", err)
	a.Error(err)
	a.Contains(err.Error(), "expected key separator")
}

func TestValidConfig(t *testing.T) {
	a := assert.New(t)
	done := prepareValidConfig(t)
	defer done()
	err := LoadConfig(&cli.Context{})
	t.Logf("found err: %s", err)
	a.NoError(err)
}

func prepare(t *testing.T) (string, func()) {
	a := assert.New(t)
	tmpDir, err := ioutil.TempDir("", "")
	a.NoError(err)
	original := userCurrent
	userCurrent = func() (*user.User, error) {
		return &user.User{HomeDir: tmpDir}, userCurrentErr
	}
	return tmpDir, func() {
		userCurrent = original
		os.RemoveAll(tmpDir)
	}
}

func prepareConfig(t *testing.T, tmpDir string, config []byte) {
	a := assert.New(t)
	file := filepath.Join(tmpDir, filename)
	dir := filepath.Dir(file)
	if st, err := os.Stat(dir); os.IsNotExist(err) || !st.IsDir() {
		a.NoError(os.MkdirAll(dir, 0700))
	}
	a.NoError(ioutil.WriteFile(file, config, 0600))
}

func prepareValidConfig(t *testing.T) func() {
	a := assert.New(t)
	tmpDir, done := prepare(t)
	prepareConfig(t, tmpDir, []byte(`
dir = '`+tmpDir+`'
ignores = [
	'\A\.git',
	'\A\.vim',
	'node_modules',
]
words = [
	'hoge',
	'fuga',
]

[mail]
from = 'hoge@example.com'
to = ['fuga@example.com']
subject = 'hoge fugao'
text = 'hoge fugafuga'
`))
	a.NoError(
		ioutil.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(""), 0600),
	)
	a.NoError(LoadConfig(&cli.Context{}))
	return done
}
