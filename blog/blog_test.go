package blog

import (
	"context"
	"fmt"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp/cmpopts"

	"encore.dev/beta/auth"
)

func TestPublishAndList(t *testing.T) {
	c := qt.New(t)
	ctx := auth.WithContext(context.Background(), "dummy", nil)
	aWhileAgo := time.Now().Add(time.Hour * -46).String()
	p := &CreateBlogPostParams{
		Title:     "title",
		Slug:      "newslug",
		Summary:   "summary",
		CreatedAt: aWhileAgo,
		Published: false,
	}

	err := CreateBlogPost(ctx, p)
	fmt.Println(err)
	c.Assert(err, qt.IsNil)

	list, err := GetBlogPosts(ctx, &GetBlogPostsParams{Limit: 100})
	fmt.Println(err)

	c.Assert(err, qt.IsNil)
	found := false
	for _, b := range list.BlogPosts {
		fmt.Println(b.Slug)
		if b.Slug == p.Slug {
			c.Assert(b.Title, qt.Equals, p.Title)
			c.Assert(b.Summary, qt.Equals, p.Summary)
			c.Assert(b.CreatedAt, qt.CmpEquals(cmpopts.EquateApproxTime(1*time.Second)), aWhileAgo)
			found = true
		}
	}
	c.Assert(found, qt.IsTrue)
}
