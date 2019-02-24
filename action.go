package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/xerrors"
	"gopkg.in/urfave/cli.v2"
)

// Action is the main logic
func Action(c *cli.Context) error {
	w, err := NewWalker()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	dir, err := targetDir()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err = filepath.Walk(dir, w.walk); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	header := resultHeader(w)
	LogBytes(c, header)
	if !c.Bool("mail") {
		return nil
	}
	subject, err := resultSubject(w)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	body := bytes.NewBuffer(header)
	_, _ = body.Write([]byte{'\n', '\n'})
	resultBody(body, w)
	if err := Mail(subject, body.Bytes()); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func targetDir() (string, error) {
	info, err := os.Lstat(Config.Dir)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return Config.Dir, nil
	}
	dir, err := filepath.EvalSymlinks(Config.Dir)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	return dir, nil
}

func resultHeader(w *Walker) []byte {
	return []byte(fmt.Sprintf(
		"dotfiles scan for prohibited words; scanned: %d files(s), "+
			"skipped: %d files(s), found: %d files(s)",
		w.scanned, w.skipped, len(w.found)))
}

func resultBody(wr io.Writer, w *Walker) {
	for _, found := range w.found {
		_, _ = wr.Write([]byte("  - "))
		_, _ = wr.Write([]byte{'\n'})
		_, _ = wr.Write([]byte(found))
	}
}

func resultSubject(w *Walker) (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	ymd := time.Now().Format("2006/01/02")
	return fmt.Sprintf("[prohibited words checker] %s %s report", host, ymd), nil
}
