// Copyright 2024-2026 Victor Dodon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package cmd

import (
	"fmt"
	"io"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

const (
	flagNameAuthGetAuthorizeLinkState = "oauth-state"
	flagNameAuthGetAuthorizeLinkOpen  = "open"
)

// authGetAuthorizeLinkCmd represents the `accounts get-list` command
var authGetAuthorizeLinkCmd = &cobra.Command{
	Use:   "get-authorize-link",
	Short: "Get OAuth2 authorize link",
	Long:  `Get OAuth2 authorize link`,
	RunE: func(cmd *cobra.Command, args []string) error {
		oauth2Cfg, err := NewOAuth2Config(cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fvState, err := cmd.Flags().GetString(flagNameAuthGetAuthorizeLinkState)
		if err != nil {
			return err
		}

		fvOpen, err := cmd.Flags().GetBool(flagNameAuthGetAuthorizeLinkOpen)
		if err != nil {
			return err
		}

		authURL := oauth2Cfg.AuthCodeURL(fvState)
		fmt.Printf("%s\n", authURL)

		if fvOpen {
			// Discard info messages from browser package printed to Stdout.
			browser.Stdout = io.Discard
			if err := browser.OpenURL(authURL); err != nil {
				cmd.SilenceUsage = true
				return err
			}
		}

		return nil
	},
}

func init() {
	authGetAuthorizeLinkCmd.Flags().String(flagNameAuthGetAuthorizeLinkState, "", "OAuth2 state param")
	authGetAuthorizeLinkCmd.Flags().Bool(flagNameAuthGetAuthorizeLinkOpen, false, "Open the authorize link in the default browser")

	AuthCmd.AddCommand(authGetAuthorizeLinkCmd)
}
