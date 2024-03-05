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
	"errors"
	"os"
)

func setupTestClientOAuth2Config(skipIfEmptyEnv bool) (oauth2Cfg *OAuth2Config, err error) {
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

	if cfg, er := MakeOAuth2Config(
		OAuth2ConfigCredentials(clientID, clientSecret),
		OAuth2ConfigRedirectURL(redirectURL),
	); er != nil {
		err = er
		return
	} else {
		oauth2Cfg = &cfg
	}
	return
}

func setupTestClient(skipIfEmptyEnv bool, oauth2Cfg *OAuth2Config) (*Client, error) {
	if oauth2Cfg == nil {
		cfg, err := setupTestClientOAuth2Config(skipIfEmptyEnv)
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

	token, err := TokenFromJSON([]byte(tokenJSON))
	if err != nil {
		return nil, err
	}

	client, err := NewClient(
		context.Background(),
		ClientOAuth2Config(*oauth2Cfg),
		ClientOAuth2InitialToken(token),
		ClientProductionEnvironment(false),
	)
	return client, err
}
