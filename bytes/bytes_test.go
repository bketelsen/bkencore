package bytes

import (
	"context"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp/cmpopts"

	"encore.dev/beta/auth"
)

func TestPublishAndList(t *testing.T) {
	c := qt.New(t)
	ctx := auth.WithContext(context.Background(), "dummy", nil)
	p := &PublishParams{
		Title:   "title",
		Summary: "summary",
		URL:     "https://example.org",
	}

	now := time.Now()
	resp, err := Publish(ctx, p)
	c.Assert(err, qt.IsNil)
	c.Assert(resp.ID, qt.Not(qt.Equals), 0)

	list, err := List(ctx, &ListParams{})
	c.Assert(err, qt.IsNil)
	found := false
	for _, b := range list.Bytes {
		if b.ID == resp.ID {
			c.Assert(b.Title, qt.Equals, p.Title)
			c.Assert(b.Summary, qt.Equals, p.Summary)
			c.Assert(b.URL, qt.Equals, p.URL)
			c.Assert(b.Created, qt.CmpEquals(cmpopts.EquateApproxTime(1*time.Second)), now)
			found = true
		}
	}
	c.Assert(found, qt.IsTrue)
}
