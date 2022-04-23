package cmd

import (
	"time"

)

type BlogPost struct {
	Slug         string
	CreatedAt    time.Time
	Published    bool
	ModifiedAt   time.Time
	Title        string
	Summary      string
	Body         string
	BodyRendered string
}

type CreateBlogPostParams struct {
	Slug      string
	Published bool
	Title     string
	Summary   string
	Body      string
}

type GetBlogPostParams struct {
	Slug string
}

type GetBlogPostsParams struct {
	Limit  int
	Offset int
}

type GetBlogPostsResponse struct {
	Count     int
	BlogPosts []*BlogPost
}
