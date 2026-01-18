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
	"github.com/printesoi/e-factura-go/pkg/oauth2"
	"github.com/spf13/cobra"
)

func NewOAuth2Config(cmd *cobra.Command) (cfg oauth2.Config, err error) {
	fvClientID, err := cmd.InheritedFlags().GetString(FlagNameOauthClientID)
	if err != nil {
		return cfg, err
	}
	fvClientSecret, err := cmd.InheritedFlags().GetString(FlagNameOauthClientSecret)
	if err != nil {
		return cfg, err
	}
	fvRedirectURL, err := cmd.Flags().GetString(FlagNameOAuthRedirectURL)
	if err != nil {
		return cfg, err
	}

	cfg, err = oauth2.MakeConfig(
		oauth2.ConfigCredentials(fvClientID, fvClientSecret),
		oauth2.ConfigRedirectURL(fvRedirectURL),
	)
	return
}

func RegisterAuthFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(FlagNameOauthClientID, "", "OAuth2 client ID")
	_ = cmd.MarkPersistentFlagRequired(FlagNameOauthClientID)
	cmd.PersistentFlags().String(FlagNameOauthClientSecret, "", "OAuth2 client secret")
	_ = cmd.MarkPersistentFlagRequired(FlagNameOauthClientSecret)
	cmd.PersistentFlags().String(FlagNameOAuthRedirectURL, "", "OAuth2 redirect URL. This needs to match one of the URLs for the OAuth2 app in SPV.")
	_ = cmd.MarkPersistentFlagRequired(FlagNameOAuthRedirectURL)
	cmd.PersistentFlags().String(FlagNameAccessToken, "", "JSON-encoded OAuth2 token. The access token should not necessarily be valid, but rather only the refresh token.")
	_ = cmd.MarkPersistentFlagRequired(FlagNameAccessToken)
}
