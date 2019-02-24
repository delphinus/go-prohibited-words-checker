package main

import (
	"html/template"
	"io"

	"golang.org/x/xerrors"
)

var tmplNames = []string{"assets/message.tmpl", "assets/mail.tmpl"}
var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFiles(tmplNames...))
}

func messageText(wr io.Writer, w *Walker) error {
	return executeTemplate(wr, "message", w)
}

func mailText(wr io.Writer, w *Walker) error {
	return executeTemplate(wr, "mail", w)
}

func executeTemplate(wr io.Writer, name string, w *Walker) error {
	if err := tmpl.ExecuteTemplate(wr, name, w.Data()); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
