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
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"encore.app/bkml/client"
)

// shortenCmd represents the shorten command
var shortenCmd = &cobra.Command{
	Use:   "shorten URL",
	Short: "Shortens a URL and returns a new, short URL that redirects to the original URL.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse the URL
		u, err := url.Parse(args[0])
		cobra.CheckErr(err)
		if u.Scheme == "" {
			return errors.New("url must be fully qualified with a scheme and a host")
		}

		// Generate the short URL
		resp, err := backend.Url.Shorten(cmd.Context(), client.UrlShortenParams{
			URL: args[0],
		})
		cobra.CheckErr(err)
		fmt.Println(resp.ShortURL)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(shortenCmd)
}
