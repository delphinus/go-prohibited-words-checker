package main

import (
	"fmt"
	"os/user"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	ignore "github.com/sabhiram/go-gitignore"
	"golang.org/x/xerrors"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/urfave/cli.v2"
)

const filename = ".config/go-prohibited-words-checker/config.toml"

// Configs is a struct of the config. Subject has the mail subject template
// that has %s as the placeholder. Ignores has regexp's to ignore files to
// search.
type Configs struct {
	Dir     string   `toml:"dir" validate:"required"`
	Ignores []string `toml:"ignores" validate:"gt=0,dive,required"`
	Mail    struct {
		From    string   `toml:"from" validate:"required"`
		Subject string   `toml:"subject" validate:"required"`
		Text    string   `toml:"text" validate:"required"`
		To      []string `toml:"to" validate:"gt=0,dive,required"`
	} `toml:"mail" validate:"required"`
	Words []string `toml:"words" validate:"gt=0,dive,required"`
}

// GitIgnore returns the matcher for .gitignore
func (c *Configs) GitIgnore() (*ignore.GitIgnore, error) {
	matcher, err := ignore.CompileIgnoreFile(filepath.Join(c.Dir, ".gitignore"))
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return matcher, nil
}

// Config is the loaded config. This is available after Before()
var Config = &Configs{}

// for testing
var userCurrent = user.Current

func configFilename() (string, error) {
	user, err := userCurrent()
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	return path.Join(user.HomeDir, filename), nil
}

// LoadConfig loads config file and returns it parsed.
func LoadConfig(*cli.Context) error {
	file, err := configFilename()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if _, err := toml.DecodeFile(file, Config); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := isValidConfig(); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func isValidConfig() (err error) {
	validate := validator.New()
	if err = validate.Struct(Config); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return
}
