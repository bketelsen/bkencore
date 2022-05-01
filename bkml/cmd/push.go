/*
Copyright © 2022 Brian Ketelsen<mail@bjk.fyi>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/spf13/cobra"

	"encore.app/bkml/client"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Upload blog posts from 'posts' directory to blog",
	Run: func(cmd *cobra.Command, args []string) {
		err := push()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}

func push() error {
	currentDirectory, err := os.Getwd() // todo
	cobra.CheckErr(err)
	postsDir := filepath.Join(currentDirectory, "posts")
	pagesDir := filepath.Join(currentDirectory, "pages")

	err = posts(postsDir)
	if err != nil {
		return err
	}
	err = pages(pagesDir)
	if err != nil {
		return err
	}
	return nil
}

func posts(postsDir string) error {
	err := filepath.Walk(postsDir, func(path string, info os.FileInfo, err error) error {
		cobra.CheckErr(err)
		// slug will be the relative path minus the extension
		cobra.CheckErr(err)
		if !info.IsDir() {
			slug := strings.Split(info.Name(), ".")

			// read the file
			f, err := os.Open(path)
			cobra.CheckErr(err)

			// create a blogpost and populate frontmatter
			var post client.BlogCreateBlogPostParams
			rest, err := frontmatter.Parse(f, &post)
			cobra.CheckErr(err)
			post.Body = string(rest)
			post.Slug = slug[0]
			if (post.CreatedAt == time.Time{}) {
				post.CreatedAt = time.Now()
			}
			if (post.ModifiedAt == time.Time{}) {
				post.ModifiedAt = post.CreatedAt
			}

			// submit to the API
			err = backend.Blog.CreateBlogPost(context.Background(), post)
			cobra.CheckErr(err)
		}

		return nil
	})
	return err
}
func pages(pagesDir string) error {
	err := filepath.Walk(pagesDir, func(path string, info os.FileInfo, err error) error {
		cobra.CheckErr(err)
		// slug will be the relative path minus the extension
		cobra.CheckErr(err)
		if !info.IsDir() {
			slug := strings.Split(info.Name(), ".")

			// read the file
			f, err := os.Open(path)
			cobra.CheckErr(err)

			// create a blogpost and populate frontmatter
			var post client.BlogCreatePageParams
			rest, err := frontmatter.Parse(f, &post)
			cobra.CheckErr(err)
			post.Body = string(rest)

			// submit to the API
			err = backend.Blog.CreatePage(context.Background(), slug[0], post)
			cobra.CheckErr(err)
		}

		return nil
	})
	return err
}
