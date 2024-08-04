// Copyright 2024 Victor Dodon
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
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "OAuth2 protocol for ANAF APIs",
}

func init() {
	authCmd.PersistentFlags().String(flagNameOauthClientID, "", "OAuth2 client ID")
	_ = authCmd.MarkPersistentFlagRequired(flagNameOauthClientID)
	authCmd.PersistentFlags().String(flagNameOauthClientSecret, "", "OAuth2 client secret")
	_ = authCmd.MarkPersistentFlagRequired(flagNameOauthClientSecret)
	authCmd.PersistentFlags().String(flagNameOAuthRedirectURL, "", "OAuth2 redirect URL. This needs to match one of the URLs for the OAuth2 app in SPV.")
	_ = authCmd.MarkPersistentFlagRequired(flagNameOAuthRedirectURL)

	rootCmd.AddCommand(authCmd)
}
