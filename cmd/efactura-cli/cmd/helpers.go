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
	"fmt"

	"github.com/spf13/cobra"

	efacturaclient "github.com/printesoi/e-factura-go/pkg/client"
	"github.com/printesoi/e-factura-go/pkg/constants"
	"github.com/printesoi/e-factura-go/pkg/efactura"
	"github.com/printesoi/e-factura-go/pkg/oauth2"
)

func newOAuth2Config(cmd *cobra.Command) (cfg oauth2.Config, err error) {
	fvClientID, err := cmd.InheritedFlags().GetString(flagNameOauthClientID)
	if err != nil {
		return cfg, err
	}
	fvClientSecret, err := cmd.InheritedFlags().GetString(flagNameOauthClientSecret)
	if err != nil {
		return cfg, err
	}
	fvRedirectURL, err := cmd.Flags().GetString(flagNameOAuthRedirectURL)
	if err != nil {
		return cfg, err
	}

	cfg, err = oauth2.MakeConfig(
		oauth2.ConfigCredentials(fvClientID, fvClientSecret),
		oauth2.ConfigRedirectURL(fvRedirectURL),
	)
	return
}

func newEfacturaClient(ctx context.Context, cmd *cobra.Command) (client *efactura.Client, err error) {
	fvProduction, err := cmd.InheritedFlags().GetBool(flagNameProduction)
	if err != nil {
		return nil, err
	}

	fvToken, err := cmd.InheritedFlags().GetString(flagNameOAuthToken)
	if err != nil {
		return nil, err
	}

	token, err := oauth2.TokenFromJSON([]byte(fvToken))
	if err != nil {
		return nil, fmt.Errorf("error loading token from JSON: %w", err)
	}

	oauth2Cfg, err := newOAuth2Config(cmd)
	if err != nil {
		return nil, fmt.Errorf("error creating oauth2 config: %w", err)
	}

	tokenSource := oauth2Cfg.TokenSource(ctx, token)
	if fvProduction {
		client, err = efactura.NewProductionClient(ctx, tokenSource)
	} else {
		client, err = efactura.NewSandboxClient(ctx, tokenSource)
	}
	return
}

func newEfacturaPublicClient(cmd *cobra.Command) (client *efactura.Client, err error) {
	publicApiClient, err := efacturaclient.NewPublicApiClient(
		efacturaclient.PublicApiClientBaseURL(constants.PublicApiBaseProd),
	)
	if err != nil {
		return nil, err
	}

	client, err = efactura.NewClient(efactura.ClientPublicApiClient(publicApiClient))
	return
}
