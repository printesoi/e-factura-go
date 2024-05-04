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
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	flagNameOAuthDeviceCode = "oauth-device-code"
)

// authExchangeCodeCmd represents the `accounts get-list` command
var authExchangeCodeCmd = &cobra.Command{
	Use:   "exchange-code",
	Short: "Exchange an OAuth device code for a token",
	Long:  `Exchange an OAuth device code for a token`,
	RunE: func(cmd *cobra.Command, args []string) error {
		oauth2Cfg, err := newOAuth2Config(cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fvCode, err := cmd.Flags().GetString(flagNameOAuthDeviceCode)
		if err != nil {
			return err
		}

		token, err := oauth2Cfg.Exchange(context.Background(), fvCode)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		tokenJSON, err := json.Marshal(token)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fmt.Println(string(tokenJSON))
		return nil
	},
}

func init() {
	authExchangeCodeCmd.Flags().String(flagNameOAuthDeviceCode, "", "OAuth2 device code to be exchanged for a token")
	_ = authExchangeCodeCmd.MarkFlagRequired(flagNameOAuthDeviceCode)

	authCmd.AddCommand(authExchangeCodeCmd)
}
