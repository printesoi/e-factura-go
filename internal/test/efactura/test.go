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

package efactura

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	xoauth2 "golang.org/x/oauth2"

	"github.com/printesoi/e-factura-go/pkg/efactura"
	"github.com/printesoi/e-factura-go/pkg/oauth2"
)

func GetTestCIF() string {
	return os.Getenv("EFACTURA_TEST_CIF")
}

// SetupTestEnvOAuth2Config creates an OAuth2Config from the environment
// suitable for creating an efactura.Client.
// If skipIfEmptyEnv if set to false and the env variables
// EFACTURA_TEST_CLIENT_ID, EFACTURA_TEST_CLIENT_SECRET,
// EFACTURA_TEST_REDIRECT_URL are not set, this method returns an error.
// If skipIfEmptyEnv is set to true and the env vars
// are not set, this method returns a nil config.
func SetupTestEnvOAuth2Config(skipIfEmptyEnv bool) (oauth2Cfg *oauth2.Config, err error) {
	clientID := os.Getenv("EFACTURA_TEST_CLIENT_ID")
	clientSecret := os.Getenv("EFACTURA_TEST_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		if skipIfEmptyEnv {
			return
		}
		err = errors.New("invalid oauth2 credentials")
		return
	}

	redirectURL := os.Getenv("EFACTURA_TEST_REDIRECT_URL")
	if redirectURL == "" {
		err = errors.New("invalid redirect URL")
		return
	}

	if cfg, er := oauth2.MakeConfig(
		oauth2.ConfigCredentials(clientID, clientSecret),
		oauth2.ConfigRedirectURL(redirectURL),
	); er != nil {
		err = er
		return
	} else {
		oauth2Cfg = &cfg
	}
	return
}

// SetupRealClient creates a real sandboxed efactura.Client (a client that
// talks to the TEST API endpoints).
func SetupRealClient(skipIfEmptyEnv bool, oauth2Cfg *oauth2.Config, onTokenChanged oauth2.TokenChangedHandler) (*efactura.Client, error) {
	if oauth2Cfg == nil {
		cfg, err := SetupTestEnvOAuth2Config(skipIfEmptyEnv)
		if err != nil {
			return nil, err
		} else if cfg == nil {
			return nil, nil
		} else {
			oauth2Cfg = cfg
		}
	}

	tokenJSON := os.Getenv("EFACTURA_TEST_INITIAL_TOKEN_JSON")
	if tokenJSON == "" {
		return nil, errors.New("Invalid initial token json")
	}

	token, err := oauth2.TokenFromJSON([]byte(tokenJSON))
	if err != nil {
		return nil, err
	}

	if onTokenChanged == nil {
		onTokenChanged = func(ctx context.Context, token *xoauth2.Token) error {
			tokenJSON, _ := json.Marshal(token)
			fmt.Printf("[E-FACTURA] token changed: %s\n", string(tokenJSON))
			return nil
		}
	}

	ctx := context.Background()
	client, err := efactura.NewSandboxClient(ctx, oauth2Cfg.TokenSourceWithChangedHandler(ctx, token, onTokenChanged))
	return client, err
}
