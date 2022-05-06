package email

import (
	"context"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

type SendParams struct {
	MessageID int64 `json:"message_id,omitempty"`
}

type SendResponse struct {
	// Sent reports whether the email was sent or not.
	Sent bool `json:"sent,omitempty"`
	// ProviderID is the unique id from the email provider for this send.
	ProviderID string `json:"provider_id,omitempty"`
}

//encore:api private method=POST path=/email/send/:messageID
func Send(ctx context.Context, messageID int64) (*SendResponse, error) {
	eb := errs.B().Meta("message_id", messageID)

	// Read the message from the database
	tx, err := sqldb.Begin(ctx)
	if err != nil {
		return nil, eb.Cause(err).Err()
	}
	defer tx.Rollback() // comitted explicitly on success
	m, err := lockMessage(ctx, tx, messageID)
	if err != nil {
		return nil, eb.Cause(err).Err()
	}

	// Ensure the user exists and is opted in.
	if status, err := isOptedIn(ctx, m.EmailAddress); err != nil {
		return nil, eb.Cause(err).Err()
	} else if !status {
		return &SendResponse{Sent: false}, nil
	}

	token, err := encodeUnsubscribeToken(m.EmailAddress, m.ID)
	if err != nil {
		return nil, err
	}

	// Add the token to the email.
	plaintext := strings.ReplaceAll(string(m.BodyText), "{{Token}}", token)
	html := strings.ReplaceAll(string(m.BodyHTML), "{{Token}}", token)

	mg := mailgun.NewMailgun(cfg.MailgunDomain, secrets.MailGunAPIKey)
	msg := mg.NewMessage(m.Sender, m.Subject, plaintext, m.EmailAddress)
	msg.SetTrackingOpens(true)
	msg.SetTrackingClicks(false)
	msg.SetHtml(html)

	// Use a different context for sending so we don't accidentally
	// send duplicate mails if the client disconnects.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mailgunID, err := sendEmail(ctx, mg, msg)
	if err != nil {
		return nil, eb.Cause(err).Code(errs.Unavailable).Err()
	}

	// Mark the email as successfully sent.
	_, err = tx.Exec(ctx, `
		UPDATE message
		SET provider_id = $2, sent_at = NOW()
		WHERE id = $1
	`, m.ID, mailgunID)
	if err != nil {
		rlog.Error("failed to store mailgun id", "err", err)
	}
	if err := tx.Commit(); err != nil {
		rlog.Error("failed to commit tx", "err", err)
	}

	return &SendResponse{Sent: true, ProviderID: mailgunID}, nil
}

// sendEmail sends an email using mailgun.
// It is a variable to allow for mocking in tests.
var sendEmail = func(ctx context.Context, mg *mailgun.MailgunImpl, msg *mailgun.Message) (providerID string, err error) {
	_, mailgunID, err := mg.Send(ctx, msg)
	return mailgunID, err
}

// message represents a single email message to be sent.
type message struct {
	ID           int64   `json:"id,omitempty"` // message id
	EmailAddress string  `json:"email_address,omitempty"`
	ProviderID   *string `json:"provider_id,omitempty"` // nil if not yet sent

	// Inlined template data
	Sender   string `json:"sender,omitempty"`
	Subject  string `json:"subject,omitempty"`
	BodyText string `json:"body_text,omitempty"`
	BodyHTML string `json:"body_html,omitempty"`
}

// lockMessage reads a message from the database and locks it for the duration of the tx.
func lockMessage(ctx context.Context, tx *sqldb.Tx, id int64) (*message, error) {
	var m message
	err := tx.QueryRow(ctx, `
		SELECT m.id, m.email_address, m.provider_id, t.sender, t.subject, t.body_text, t.body_html
		FROM message m
		INNER JOIN template t ON (t.id = m.template_id)
		WHERE m.id = $1
		FOR UPDATE
	`, id).Scan(&m.ID, &m.EmailAddress, &m.ProviderID, &m.Sender, &m.Subject, &m.BodyText, &m.BodyHTML)
	return &m, err
}
