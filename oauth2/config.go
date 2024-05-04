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

package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/printesoi/e-factura-go/errors"
	xoauth2 "golang.org/x/oauth2"
)

const (
	// AuthURL is the ANAF authorize URL for the OAuth2 protocol
	AuthURL = "https://logincert.anaf.ro/anaf-oauth2/v1/authorize"
	// TokenURL is the ANAF token URL for the OAuth2 protocol
	TokenURL = "https://logincert.anaf.ro/anaf-oauth2/v1/token"
	// RevokeTokenURL is the ANAF token revocation URL for the OAuth2 protocol
	RevokeTokenURL = "https://logincert.anaf.ro/anaf-oauth2/v1/revoke"
)

var (
	// Endpoint is the default endpoint for ANAF OAuth2 protocol
	Endpoint = xoauth2.Endpoint{
		AuthURL:   AuthURL,
		TokenURL:  TokenURL,
		AuthStyle: xoauth2.AuthStyleInHeader,
	}
)

// Config is a wrapper over the golang.org/x/oauth2.Config.
type Config struct {
	xoauth2.Config
}

// ConfigOption allows gradually modifying a Config
type ConfigOption func(*Config)

// ConfigCredentials set client ID and client secret
func ConfigCredentials(clientID, clientSecret string) ConfigOption {
	return func(c *Config) {
		c.ClientID, c.ClientSecret = clientID, clientSecret
	}
}

// ConfigRedirectURL set the redirect URL
func ConfigRedirectURL(redirectURL string) ConfigOption {
	return func(c *Config) {
		c.RedirectURL = redirectURL
	}
}

// ConfigEndpoint sets the auth endpoint for the config. This should only
// be used with debugging/testing auth requests.
func ConfigEndpoint(endpoint xoauth2.Endpoint) ConfigOption {
	return func(c *Config) {
		c.Config.Endpoint = endpoint
	}
}

// MakeConfig creates a Config using provided options. At least
// ConfigCredentials must be provided, otherwise
// ErrInvoiceOAuth2Credentials will be returned. If an invalid endpoint if
// provided using ConfigEndpoint, then ErrInvalidOAuth2Endpoint is
// returned.
func MakeConfig(opts ...ConfigOption) (cfg Config, err error) {
	cfg.Endpoint = Endpoint
	for _, opt := range opts {
		opt(&cfg)
	}
	if !cfg.validCredentials() {
		err = errors.ErrInvalidOAuth2Credentials
		return
	}
	if !cfg.validEndpoint() {
		err = errors.ErrInvalidOAuth2Endpoint
		return
	}
	if cfg.RedirectURL == "" {
		err = errors.ErrInvalidOAuth2RedirectURL
		return
	}
	return
}

func (c Config) validCredentials() bool {
	return c.ClientID != "" && c.ClientSecret != ""
}

func (c Config) validEndpoint() bool {
	return c.Endpoint.AuthURL != "" && c.Endpoint.TokenURL != ""
}

// Valid returns true if the config is valid (ie. is non-nil, has non-empty
// credentials, non-empty endpoint, and non-empty redirect URL).
func (c *Config) Valid() bool {
	return c != nil && c.validCredentials() && c.validEndpoint() && c.RedirectURL != ""
}

// AuthCodeURL generates the code authorization URL.
func (c Config) AuthCodeURL(state string) string {
	return c.Config.AuthCodeURL(state,
		xoauth2.SetAuthURLParam("token_content_type", "jwt"))
}

// Exchange converts an authorization code into a token.
func (c Config) Exchange(ctx context.Context, code string) (*xoauth2.Token, error) {
	return c.Config.Exchange(ctx, code,
		xoauth2.SetAuthURLParam("token_content_type", "jwt"))
}

type tokenJSON struct {
	AccessToken  string     `json:"access_token"`
	TokenType    string     `json:"token_type"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int64      `json:"expires_in,omitempty"`
	Expiry       *time.Time `json:"expiry,omitempty"`
}

// timeNow is time.Now but pulled out as a variable for tests.
var timeNow = time.Now

func (tj *tokenJSON) expiry() (t time.Time) {
	if e := tj.Expiry; e != nil {
		return *e
	}
	if ei := tj.ExpiresIn; ei != 0 {
		return timeNow().Add(time.Duration(ei) * time.Second)
	}
	return
}

// TokenFromJSON is a convenience function that parses a xoauth2.Token from a
// JSON encoded value. This can parse both a JSON from the ANAF OAuth2 provider
// (with an expires_in field) or a JSON encoded xoauth2.Token (with an expiry
// field).
func TokenFromJSON(jsonData []byte) (token *xoauth2.Token, err error) {
	var tj tokenJSON
	if err = json.Unmarshal(jsonData, &tj); err != nil {
		err = fmt.Errorf("efactura.oauth2: cannot parse json: %v", err)
		return
	}

	t := xoauth2.Token{
		AccessToken:  tj.AccessToken,
		TokenType:    tj.TokenType,
		RefreshToken: tj.RefreshToken,
		Expiry:       tj.expiry(),
	}
	if t.Type() != "Bearer" || t.AccessToken == "" || t.RefreshToken == "" {
		err = fmt.Errorf("efactura.oauth2: malformed or incomplete token")
		return
	}
	return &t, nil
}
