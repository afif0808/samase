package samasemailservice

import "context"

type SendEmailFunc func(ctx context.Context, dest []string, subject string, body string, attachments ...string) error
