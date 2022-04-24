package email

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/mailgun/mailgun-go/v4"

	"encore.dev/beta/auth"
)

func TestEndToEnd(t *testing.T) {
	c := qt.New(t)
	sent := useSendMock(c)
	ctx := auth.WithContext(context.Background(), "user-id", nil)

	// Create an email template
	tmplID := c.Name()
	err := CreateTemplate(ctx, tmplID, &CreateTemplateParams{
		Sender:   "sender@example.org",
		Subject:  "subject",
		BodyText: "text body",
		BodyHTML: "html body",
	})
	c.Assert(err, qt.IsNil)

	// Schedule an email to be sent to two people
	scheduled, err := Schedule(ctx, &ScheduleParams{
		TemplateID:     tmplID,
		EmailAddresses: []string{"john@example.org", "jane@example.org"},
	})
	c.Assert(err, qt.IsNil)
	c.Assert(scheduled.MessageIDs, qt.HasLen, 2)

	// Send the emails
	for _, id := range scheduled.MessageIDs {
		resp, err := Send(ctx, id)
		c.Assert(err, qt.IsNil)
		c.Assert(resp, qt.DeepEquals, &SendResponse{Sent: true, ProviderID: "msg-id"})
	}
	c.Assert(*sent, qt.HasLen, len(scheduled.MessageIDs))
}

func useSendMock(c *qt.C) *[]*mailgun.Message {
	var msgs []*mailgun.Message
	orig := sendEmail
	sendEmail = func(ctx context.Context, mg *mailgun.MailgunImpl, msg *mailgun.Message) (providerID string, err error) {
		msgs = append(msgs, msg)
		return "msg-id", nil
	}
	c.Cleanup(func() {
		sendEmail = orig
	})
	return &msgs
}
