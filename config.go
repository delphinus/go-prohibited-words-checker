package main

import (
	"os/user"
	"path"

	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
	"gopkg.in/urfave/cli.v2"
)

const filename = ".config/go-prohibited-words-checker/config.toml"

// Config is a struct of the config. Subject has the mail subject template
// that has %s as the placeholder.
type Configs struct {
	Words []string `toml:"words"`
	Mail  struct {
		To      []string `toml:"to"`
		From    string   `toml:"from"`
		Subject string   `toml:"subject"`
		Text    string   `toml:"text"`
	} `toml:"mail"`
}

// Config is the loaded config. This is available after Before()
var Config *Configs

func configFilename() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", xerrors.New("not found the current user")
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
	return nil
}
