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

	"encore.app/blogsync/client"
	"github.com/adrg/frontmatter"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := push()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func push() error {
	var base = client.BaseURL("http://127.0.0.1:4000") // todo
	cl, err := client.New(base, client.WithAuthToken("yourmom"))
	cobra.CheckErr(err)

	currentDirectory, err := os.Getwd() // todo
	cobra.CheckErr(err)
	postsDir := filepath.Join(currentDirectory, "posts")

	filepath.Walk(postsDir, func(path string, info os.FileInfo, err error) error {
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

			// submit to the API
			err = cl.Blog.CreateBlogPost(context.Background(), post)
			cobra.CheckErr(err)
		}

		return nil
	})

	return nil
}
