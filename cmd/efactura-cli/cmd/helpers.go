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
	xoauth2 "golang.org/x/oauth2"

	efacturaclient "github.com/printesoi/e-factura-go/pkg/client"
	"github.com/printesoi/e-factura-go/pkg/constants"
	"github.com/printesoi/e-factura-go/pkg/efactura"
	"github.com/printesoi/e-factura-go/pkg/oauth2"
)

var (
	clientCtxKey = struct{}{}
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

func getContextClient(ctx context.Context) (client *efactura.Client) {
	ctxValue := ctx.Value(clientCtxKey)
	if ctxValue == nil {
		return nil
	}
	if client, ok := ctxValue.(*efactura.Client); ok {
		return client
	}
	return nil
}

func contextWithClient(ctx context.Context, client *efactura.Client) context.Context {
	return context.WithValue(ctx, clientCtxKey, client)
}

func newEfacturaClient(ctx context.Context, cmd *cobra.Command) (client *efactura.Client, err error) {
	if client = getContextClient(ctx); client != nil {
		return
	}

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

	onTokenChanged := func(ctx context.Context, token *xoauth2.Token) error {
		tokenJSON, _ := json.Marshal(token)
		fmt.Printf("[E-FACTURA] token changed: %s\n", string(tokenJSON))
		return nil
	}
	tokenSource := oauth2Cfg.TokenSourceWithChangedHandler(ctx, token, onTokenChanged)
	if fvProduction {
		client, err = efactura.NewProductionClient(ctx, tokenSource)
	} else {
		client, err = efactura.NewSandboxClient(ctx, tokenSource)
	}
	return
}

func newEfacturaPublicClient(_ *cobra.Command) (client *efactura.Client, err error) {
	publicApiClient, err := efacturaclient.NewPublicApiClient(
		efacturaclient.PublicApiClientBaseURL(constants.PublicApiBaseProd),
	)
	if err != nil {
		return nil, err
	}

	client, err = efactura.NewClient(efactura.ClientPublicApiClient(publicApiClient))
	return
}
