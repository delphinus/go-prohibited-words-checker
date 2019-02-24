package main

import (
	"net/smtp"

	"github.com/jordan-wright/email"
	"golang.org/x/xerrors"
)

const hostport = "smtp.gmail.com:587"
const host = "smtp.gmail.com"

// Mail is a func to mail
func Mail(subject string, body []byte) error {
	e := &email.Email{
		To:      Config.Mail.To,
		From:    Config.Mail.From,
		Subject: subject,
		Text:    body,
	}
	if err := e.Send(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", Config.Mail.From, Config.Mail.Password, host),
	); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
