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
	"os"

	"github.com/spf13/cobra"

	"encore.app/bkml/client"
)

var (
	// backend is the client to communicate with the backend.
	backend *client.Client

	// envName is the backend env name to communicate with.
	// "local" means local develoment.
	envName string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bkml",
	Short: "Commands to administer content",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		base := client.Local
		if envName != "local" {
			base = client.Environment(envName)
		}
		token := os.Getenv("AUTH_PASSWORD")
		if token == "" {
			return errors.New("no AUTH_PASSWORD set")
		}
		fmt.Printf("Using %s environment\n", envName)
		backend, err = client.New(base, client.WithAuthToken(token))
		return err
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&envName, "env", "e", "staging", "environment to connect to ('local',  'staging')")
}
