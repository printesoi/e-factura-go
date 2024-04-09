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
	"math"
	"time"

	"golang.org/x/oauth2"
)

var (
	// Endpoint is the default endpoint for ANAF OAuth2 protocol
	Endpoint = oauth2.Endpoint{
		AuthURL:   "https://logincert.anaf.ro/anaf-oauth2/v1/authorize",
		TokenURL:  "https://logincert.anaf.ro/anaf-oauth2/v1/token",
		AuthStyle: oauth2.AuthStyleInHeader,
	}

	ErrInvalidOAuth2Credentials = errors.New("Invalid OAuth2 credentials")
	ErrInvalidOAuth2Endpoint    = errors.New("Invalid OAuth2 endpoint")
)

type OAuth2Config struct {
	oauth2.Config
}

// OAuth2ConfigOption allows gradually modifying an OAuth2Config
type OAuth2ConfigOption func(*OAuth2Config)

// OAuth2ConfigCredentials set client ID and client secret
func OAuth2ConfigCredentials(clientID, clientSecret string) OAuth2ConfigOption {
	return func(c *OAuth2Config) {
		c.ClientID, c.ClientSecret = clientID, clientSecret
	}
}

// OAuth2ConfigRedirectURL set the redirect URL
func OAuth2ConfigRedirectURL(redirectURL string) OAuth2ConfigOption {
	return func(c *OAuth2Config) {
		c.RedirectURL = redirectURL
	}
}

// OAuth2ConfigEndpoint sets the auth endpoint for the config. This should only
// be used with debugging/testing auth requests.
func OAuth2ConfigEndpoint(endpoint oauth2.Endpoint) OAuth2ConfigOption {
	return func(c *OAuth2Config) {
		c.Config.Endpoint = endpoint
	}
}

// MakeOAuth2Config creates a OAuth2Config using provided options. At least
// OAuth2ConfigCredentials must be provided, otherwise
// ErrInvoiceOAuth2Credentials will be returned. If an invalid endpoint if
// provided usingaOAuth2ConfigEndpoint, then ErrInvalidOAuth2Endpoint is
// returned.
func MakeOAuth2Config(opts ...OAuth2ConfigOption) (cfg OAuth2Config, err error) {
	cfg.Endpoint = Endpoint
	for _, opt := range opts {
		opt(&cfg)
	}
	if !cfg.validCredentials() {
		err = ErrInvalidOAuth2Credentials
	}
	if !cfg.validEndpoint() {
		err = ErrInvalidOAuth2Endpoint
	}
	return
}

func (c OAuth2Config) validCredentials() bool {
	return c.ClientID != "" && c.ClientSecret != ""
}

func (c OAuth2Config) validEndpoint() bool {
	return c.Endpoint.AuthURL != "" && c.Endpoint.TokenURL != ""
}

// Valid returns true if the config is valid (ie. is non-nil, has non-empty
// credentials, non-empty endpoint, and non-empty redirect URL).
func (c *OAuth2Config) Valid() bool {
	return c != nil && c.validCredentials() && c.validEndpoint() && c.RedirectURL != ""
}

// AuthCodeURL generates the code authorization URL.
func (c OAuth2Config) AuthCodeURL(state string) string {
	return c.Config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("token_content_type", "jwt"))
}

// Exchange converts an authorization code into a token.
func (c OAuth2Config) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.Config.Exchange(ctx, code,
		oauth2.SetAuthURLParam("token_content_type", "jwt"))
}

// tokenJSON is the struct representing the HTTP response from OAuth2
// providers returning a token or error in JSON form.
// https://datatracker.ietf.org/doc/html/rfc6749#section-5.1
type tokenJSON struct {
	AccessToken  string         `json:"access_token"`
	TokenType    string         `json:"token_type"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    expirationTime `json:"expires_in"`
}

func (e *tokenJSON) expiry() (t time.Time) {
	if v := e.ExpiresIn; v != 0 {
		return timeNow().Add(time.Duration(v) * time.Second)
	}
	return
}

type expirationTime int32

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	if i > math.MaxInt32 {
		i = math.MaxInt32
	}
	*e = expirationTime(i)
	return nil
}

// TokenFromJSON is a convenience method that parses an oauth2.Token from a
// JSON encoded value.
func TokenFromJSON(body []byte) (token *oauth2.Token, err error) {
	var tj tokenJSON
	if err = json.Unmarshal(body, &tj); err != nil {
		err = fmt.Errorf("oauth2: cannot parse json: %v", err)
		return
	}

	token = &oauth2.Token{
		AccessToken:  tj.AccessToken,
		TokenType:    tj.TokenType,
		RefreshToken: tj.RefreshToken,
		Expiry:       tj.expiry(),
	}
	return
}
