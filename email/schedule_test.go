package email

import (
	"context"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	"encore.dev/beta/auth"
)

func TestQueryDueEmails(t *testing.T) {
	c := qt.New(t)
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

	// Schedule two messages to be sent now, and one to be sent in a minute.
	now, future := time.Now(), time.Now().Add(time.Minute)
	scheduled, err := Schedule(ctx, &ScheduleParams{
		TemplateID:     tmplID,
		EmailAddresses: []string{"john@example.org", "jane@example.org"},
		SendAt:         &now,
	})
	c.Assert(err, qt.IsNil)
	_, err = Schedule(ctx, &ScheduleParams{
		TemplateID:     tmplID,
		EmailAddresses: []string{"doc.brown@example.org"},
		SendAt:         &future,
	})
	c.Assert(err, qt.IsNil)

	// The first two emails should be due now, but not the third.
	due, err := queryDueEmails(ctx, now)
	c.Assert(err, qt.IsNil)
	c.Assert(due, qt.DeepEquals, scheduled.MessageIDs)
}
