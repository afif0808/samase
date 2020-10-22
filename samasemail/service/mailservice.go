package samasemailservice

import (
	"context"

	"gopkg.in/gomail.v2"
)

func SendEmail(dialer *gomail.Dialer, mailer *gomail.Message) SendEmailFunc {
	return func(ctx context.Context, dest []string, subject string, body string, attachments ...string) error {
		mailer.SetHeader("To", dest...)
		mailer.SetHeader("Subject", subject)
		mailer.SetBody("text/html", body)
		for _, v := range attachments {
			mailer.Attach(v)
		}
		return dialer.DialAndSend(mailer)
	}
}
