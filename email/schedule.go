package email

import (
	"context"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"encore.dev/cron"
	"encore.dev/storage/sqldb"
)

type ScheduleParams struct {
	// TemplateID is the email template to use.
	TemplateID string

	// EmailAddresses are the email addresses to send to.
	EmailAddresses []string

	// SendAt is the time to send the email.
	// If nil it defaults to the current time.
	SendAt *time.Time
}

type ScheduleResponse struct {
	// MessageIDs are the ids for the messages that were scheduled.
	MessageIDs []int64
}

// Schedule schedules emails to be sent.
func Schedule(ctx context.Context, p *ScheduleParams) (*ScheduleResponse, error) {
	// Ensure the emails all exist in the user database.
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "user" (email_address, optin, optin_changed)
		VALUES (unnest($1::text[]), true, NOW())
		ON CONFLICT (email_address) DO NOTHING
	`, p.EmailAddresses)
	if err != nil {
		return nil, err
	}

	// Schedule the email messages to be sent.
	sendAt := time.Now()
	if p.SendAt != nil {
		sendAt = *p.SendAt
	}

	rows, err := sqldb.Query(ctx, `
		INSERT INTO "message" (email_address, template_id, scheduled_at)
		VALUES (unnest($1::text[]), $2, $3)
		RETURNING id
	`, p.EmailAddresses, p.TemplateID, sendAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &ScheduleResponse{MessageIDs: ids}, nil
}

type ScheduleAllParams struct {
	// TemplateID is the email template to use.
	TemplateID string

	// SendAt is the time to send the email.
	// If nil it defaults to the current time.
	SendAt *time.Time
}

// ScheduleAll schedules emails to be sent to all subscribers.
func ScheduleAll(ctx context.Context, p *ScheduleAllParams) (*ScheduleResponse, error) {
	// Schedule the email messages to be sent.
	sendAt := time.Now()
	if p.SendAt != nil {
		sendAt = *p.SendAt
	}

	rows, err := sqldb.Query(ctx, `
		INSERT INTO "message" (email_address, template_id, scheduled_at)
		SELECT email, $2, $3
			FROM "user"
			WHERE optin
		RETURNING id
	`, p.TemplateID, sendAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &ScheduleResponse{MessageIDs: ids}, nil
}

type CreateTemplateParams struct {
	Sender   string // sender email
	Subject  string // subject line to use
	BodyText string // plaintext body
	BodyHTML string // html body
}

// CreateTemplate creates an email template.
// If the template with that id already exists it is updated.
//encore:api auth method=PUT path=/email/templates/:id
func CreateTemplate(ctx context.Context, id string, p *CreateTemplateParams) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "template" (id, sender, subject, body_text, body_html)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
		sender = $2, subject = $3, body_text = $4, body_html = $5, updated_at = NOW()
	`, id, p.Sender, p.Subject, p.BodyText, p.BodyHTML)
	return err
}

type SendDueEmailsResponse struct {
	NumSent int // number of emails successfully sent
}

// SendDueEmails sends emails that are due to for delivery.
//encore:api private method=POST path=/email/send-due
func SendDueEmails(ctx context.Context) (*SendDueEmailsResponse, error) {
	ids, err := queryDueEmails(ctx, time.Now())
	if err != nil && len(ids) == 0 {
		return nil, err
	}

	g, ctx := errgroup.WithContext(ctx)
	var successes int64
	for _, id := range ids {
		id := id // capture for closure
		g.Go(func() error {
			_, err := Send(ctx, id)
			if err == nil {
				atomic.AddInt64(&successes, 1)
			}
			return err
		})
	}

	err = g.Wait()
	if successes > 0 {
		// If we successfully sent some emails always treat it as successful.
		err = nil
	}
	return &SendDueEmailsResponse{NumSent: int(successes)}, err
}

// queryDueEmails reports the message ids that are due to be sent at time t.
func queryDueEmails(ctx context.Context, t time.Time) (ids []int64, err error) {
	rows, err := sqldb.Query(ctx, `
		SELECT id
		FROM message
		WHERE sent_at IS NULL AND scheduled_at <= $1
		LIMIT 100
	`, t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			break
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// Send emails due to be delivered every minute.
var _ = cron.NewJob("send-due-emails", cron.JobConfig{
	Title:    "Send due emails",
	Every:    1 * cron.Minute,
	Endpoint: SendDueEmails,
})
