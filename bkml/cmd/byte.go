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

func init() {
	var title, desc string

	// byteCmd represents the shorten command
	var byteCmd = &cobra.Command{
		Use:   "byte URL --title=TITLE --desc=DESC",
		Short: "Publish a quick byte, a short post that links to interesting articles etc",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse the URL
			u, err := url.Parse(args[0])
			cobra.CheckErr(err)
			if u.Scheme == "" {
				return errors.New("url must be fully qualified with a scheme and a host")
			}

			if title == "" {
				return errors.New("--title must not be empty")
			}

			// Generate the short URL
			resp, err := backend.Bytes.Publish(cmd.Context(), client.BytesPublishParams{
				Title:   title,
				Summary: desc,
				URL:     args[0],
			})
			cobra.CheckErr(err)
			fmt.Printf("Successfully published byte with id %v\n", resp.ID)
			return nil
		},
	}

	byteCmd.Flags().StringVar(&title, "title", "", "The byte's title (required)")
	byteCmd.Flags().StringVar(&desc, "desc", "", "The byte's description (optional)")
	rootCmd.AddCommand(byteCmd)
}
