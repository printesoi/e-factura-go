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

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/printesoi/e-factura-go/constants"
	api_helpers "github.com/printesoi/e-factura-go/internal/helpers/api"
	"github.com/printesoi/e-factura-go/oauth2"
	"github.com/stretchr/testify/assert"
	xoauth2 "golang.org/x/oauth2"
)

// setupTestOAuth2Config sets up a test HTTP server along with a OAuth2Config
// that is configured to talk to that test server.
func setupTestOAuth2Config(clientID, clientSecret string) (oauth2Cfg oauth2.Config, mux *http.ServeMux, serverURL string, teardown func(), err error) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	authServer := httptest.NewServer(mux)

	serverURL = authServer.URL
	authorizeURL, err := api_helpers.BuildParseURL(serverURL, "/authorize", nil)
	if err != nil {
		return
	}
	tokenURL, err := api_helpers.BuildParseURL(serverURL, "/token", nil)
	if err != nil {
		return
	}
	redirectURL, err := api_helpers.BuildParseURL(serverURL, "/redirect", nil)
	if err != nil {
		return
	}
	oauth2Cfg, err = oauth2.MakeConfig(
		oauth2.ConfigCredentials(clientID, clientSecret),
		oauth2.ConfigRedirectURL(redirectURL),
		oauth2.ConfigEndpoint(xoauth2.Endpoint{
			AuthURL:   authorizeURL,
			TokenURL:  tokenURL,
			AuthStyle: xoauth2.AuthStyleInHeader,
		}),
	)
	if err != nil {
		return
	}

	teardown = func() {
		authServer.Close()
	}
	return
}

// setupTestApiClient sets up a test HTTP server along with a Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setupTestApiClient(oauth2Cfg oauth2.Config, token *xoauth2.Token, basePath string) (client *ApiClient, mux *http.ServeMux, serverURL string, teardown func(), err error) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/", mux)

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)
	teardown = func() {
		server.Close()
	}

	serverURL = server.URL
	baseURL, err := api_helpers.BuildParseURL(serverURL, basePath, nil)
	if err != nil {
		return
	}
	ctx := context.Background()
	client, err = NewApiClient(
		ApiClientOAuth2TokenSource(oauth2Cfg.TokenSource(ctx, token)),
		ApiClientSandboxEnvironment(true),
		ApiClientBaseURL(baseURL),
		ApiClientContext(ctx),
	)
	return
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func TestApiClientAuth(t *testing.T) {
	assert := assert.New(t)

	baseAccessToken := fmt.Sprintf("testAccessToken:%s", t.Name())
	baseRefreshToken := fmt.Sprintf("testRefreshToken:%s", t.Name())

	buildAccessToken := func(seq int64) string {
		return fmt.Sprintf("%s:%d", baseAccessToken, seq)
	}
	buildRefreshToken := func(seq int64) string {
		return fmt.Sprintf("%s:%d", baseRefreshToken, seq)
	}

	clientID, clientSecret := "test_client_id", "test_client_secret"
	code := "1234567890"

	oauth2Cfg, authMux, _, authTeardown, err := setupTestOAuth2Config(clientID, clientSecret)
	if authTeardown != nil {
		defer authTeardown()
	}
	if !assert.NoError(err) {
		return
	}

	authSeq := int64(0)
	authMux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		generateToken := func(seq int64) {
			type tokenT struct {
				AccessToken  string `json:"access_token"`
				TokenType    string `json:"token_type,omitempty"`
				RefreshToken string `json:"refresh_token,omitempty"`
				ExpiresIn    int    `json:"expires_in,omitempty"`
			}

			var token tokenT
			token.TokenType = "bearer"
			token.AccessToken = buildAccessToken(seq)
			token.RefreshToken = buildRefreshToken(seq)
			// 12 seconds = 2 seconds of validity (oauth2.defaultExpiryDelta = 10s)
			token.ExpiresIn = 12
			w.Header().Add("Content-Type", api_helpers.MediaTypeApplicationJSON)

			tokenBytes, err := json.Marshal(token)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(tokenBytes)
		}

		grantType := r.Form.Get("grant_type")
		switch grantType {
		case "authorization_code":
			// Ensure the token request has correct params
			assert.Equal(code, r.Form.Get("code"))
			assert.Equal("jwt", r.Form.Get("token_content_type"))
			generateToken(authSeq)

		case "refresh_token":
			// Ensure the token request has correct params
			assert.Equal(buildRefreshToken(authSeq-1), r.Form.Get("refresh_token"))
			generateToken(authSeq)

		default:
			t.Errorf("Unknown grant_type '%q'", grantType)
			return
		}

		authSeq++
	})

	token, err := oauth2Cfg.Exchange(context.Background(), code)
	if !assert.NoError(err) {
		return
	}

	basePath := constants.ApiBasePathSandbox
	client, clientMux, serverURL, clientTeardown, err := setupTestApiClient(oauth2Cfg, token, basePath)
	if clientTeardown != nil {
		defer clientTeardown()
	}
	if !assert.NoError(err) {
		return
	}
	_, _, _ = client, clientMux, serverURL

	assert.Equal(int64(1), authSeq, "must have generated 1 token")

	clientMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Fail("handler not properly registered")
	})

	// All requests should have the Authorization header properly set
	testAuthPath := func(path string) {
		path, _ = url.JoinPath("/", basePath, path)

		var pathHandled atomic.Int32
		clientMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			_ = pathHandled.Add(1)

			val := fmt.Sprintf("Bearer %s", buildAccessToken(authSeq-1))
			testHeader(t, r, "Authorization", val)
		})

		req, err := client.NewRequest(context.Background(), http.MethodGet, path, nil, nil)
		if !assert.NoError(err) {
			return
		}
		_, err = client.Do(req)
		if !assert.NoError(err) {
			return
		}

		assert.True(pathHandled.Load() != 0, "client should build correct path")
	}

	testAuthPath("/test_01")

	assert.Equal(int64(1), authSeq, "a call with a valid token must not trigger token refresh")
	<-time.After(2 * time.Second)

	testAuthPath("/test_02")

	assert.Equal(int64(2), authSeq, "after the token expired, it must be refreshed")
}
