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

package etransport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	xoauth2 "golang.org/x/oauth2"

	"github.com/printesoi/e-factura-go/pkg/etransport"
	"github.com/printesoi/e-factura-go/pkg/oauth2"
)

func GetTestCIF() string {
	return os.Getenv("ETRANSPORT_TEST_CIF")
}

// SetupTestEnvOAuth2Config creates a OAuth2Config from the environment
// suitable for creating an etransport.Client.
// If skipIfEmptyEnv if set to false and the env variables
// ETRANSPORT_TEST_CLIENT_ID, ETRANSPORT_TEST_CLIENT_SECRET,
// ETRANSPORT_TEST_REDIRECT_URL are not set, this method returns an error.
// If skipIfEmptyEnv is set to true and the env vars
// are not set, this method returns a nil config.
func SetupTestEnvOAuth2Config(skipIfEmptyEnv bool) (oauth2Cfg *oauth2.Config, err error) {
	clientID := os.Getenv("ETRANSPORT_TEST_CLIENT_ID")
	clientSecret := os.Getenv("ETRANSPORT_TEST_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		if skipIfEmptyEnv {
			return
		}
		err = errors.New("invalid oauth2 credentials")
		return
	}

	redirectURL := os.Getenv("ETRANSPORT_TEST_REDIRECT_URL")
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

// SetupRealClient creates a real sandboxed etransport.Client (a client that
// talks to the ANAF TEST API endpoints).
func SetupRealClient(skipIfEmptyEnv bool, oauth2Cfg *oauth2.Config, onTokenChanged oauth2.TokenChangedHandler) (*etransport.Client, error) {
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

	tokenJSON := os.Getenv("ETRANSPORT_TEST_INITIAL_TOKEN_JSON")
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
			fmt.Printf("[E-TRANSPORT] token changed: %s\n", string(tokenJSON))
			return nil
		}
	}

	ctx := context.Background()
	client, err := etransport.NewSandboxClient(ctx, oauth2Cfg.TokenSourceWithChangedHandler(ctx, token, onTokenChanged))
	return client, err
}
