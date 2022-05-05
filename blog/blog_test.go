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

	cat := &Category{
		Category: "foo",
		Summary:  "bar",
	}
	err := CreateCategory(ctx, cat)

	c.Assert(err, qt.IsNil)
	aWhileAgo := time.Now().Add(time.Hour * -46)
	isoAgo := aWhileAgo.Format(time.RFC3339)
	p := &CreateBlogPostParams{
		Title:      "title",
		Slug:       "newslug",
		Summary:    "summary",
		CreatedAt:  isoAgo,
		ModifiedAt: isoAgo,
		Published:  false,
		Category:   "foo",
	}

	err = CreateBlogPost(ctx, p)
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
