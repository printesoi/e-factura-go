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
	"context"
	"encoding/json"
	"fmt"

	"github.com/printesoi/e-factura-go/pkg/oauth2"
	"github.com/spf13/cobra"
)

const (
	flagNameAuthAccessToken = "access-token"
)

// authRefreshTokenCmd represents the `accounts get-list` command
var authRefreshTokenCmd = &cobra.Command{
	Use:   "refresh-token",
	Short: "Refresh an access token",
	Long:  `Refresh an access token`,
	RunE: func(cmd *cobra.Command, args []string) error {
		oauth2Cfg, err := NewOAuth2Config(cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fvAccessToken, err := cmd.Flags().GetString(flagNameAuthAccessToken)
		if err != nil {
			return err
		}

		accessToken, err := oauth2.TokenFromJSON([]byte(fvAccessToken))
		if err != nil {
			return err
		}

		refresher := oauth2Cfg.TokenRefresher(context.Background(), accessToken, nil)
		newToken, err := refresher.Token()
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		tokenJSON, err := json.Marshal(newToken)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fmt.Println(string(tokenJSON))
		return nil
	},
}

func init() {
	authRefreshTokenCmd.Flags().String(flagNameAuthAccessToken, "", "Access token")
	_ = authRefreshTokenCmd.MarkFlagRequired(flagNameAuthAccessToken)

	AuthCmd.AddCommand(authRefreshTokenCmd)
}
