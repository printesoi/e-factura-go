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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/printesoi/e-factura-go/oauth2"
	"github.com/stretchr/testify/assert"
	xoauth2 "golang.org/x/oauth2"
)

// setupTestEnvOAuth2Config creates a OAuth2Config from the environment.
// If skipIfEmptyEnv if set to false and the env variables
// EFACTURA_TEST_CLIENT_ID, EFACTURA_TEST_CLIENT_SECRET,
// EFACTURA_TEST_REDIRECT_URL are not set, this method returns an error.
// If skipIfEmptyEnv is set to true and the env vars
// are not set, this method returns a nil config.
func setupTestEnvOAuth2Config(skipIfEmptyEnv bool) (oauth2Cfg *oauth2.Config, err error) {
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

func getTestCIF() string {
	return os.Getenv("EFACTURA_TEST_CIF")
}

// setupRealClient creates a real sandboxed Client (a client that talks to the
// ANAF TEST APIs).
func setupRealClient(skipIfEmptyEnv bool, oauth2Cfg *oauth2.Config) (*Client, error) {
	if oauth2Cfg == nil {
		cfg, err := setupTestEnvOAuth2Config(skipIfEmptyEnv)
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

	sandbox := true
	if os.Getenv("EFACTURA_TEST_PRODUCTION") == getTestCIF() {
		sandbox = false
	}

	onTokenChanged := func(ctx context.Context, token *xoauth2.Token) error {
		tokenJSON, _ := json.Marshal(token)
		fmt.Printf("[E-FACTURA] token changed: %s\n", string(tokenJSON))
		return nil
	}

	ctx := context.Background()
	client, err := NewClient(ctx,
		ClientOAuth2TokenSource(oauth2Cfg.TokenSourceWithChangedHandler(ctx, token, onTokenChanged)),
		ClientSandboxEnvironment(sandbox),
	)
	return client, err
}

// setupTestOAuth2Config sets up a test HTTP server along with a OAuth2Config
// that is configured to talk to that test server.
func setupTestOAuth2Config(clientID, clientSecret string) (oauth2Cfg oauth2.Config, mux *http.ServeMux, serverURL string, teardown func(), err error) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	authServer := httptest.NewServer(mux)

	serverURL = authServer.URL
	authorizeURL, err := buildParseURL(serverURL, "/authorize", nil)
	if err != nil {
		return
	}
	tokenURL, err := buildParseURL(serverURL, "/token", nil)
	if err != nil {
		return
	}
	redirectURL, err := buildParseURL(serverURL, "/redirect", nil)
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

// setupTestClient sets up a test HTTP server along with a Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setupTestClient(token *xoauth2.Token) (client *Client, mux *http.ServeMux, serverURL string, teardown func(), err error) {
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
	client, err = NewClient(
		context.Background(),
		ClientOAuth2TokenSource(xoauth2.StaticTokenSource(token)),
		ClientSandboxEnvironment(true),
		ClientBaseURL(serverURL+apiBasePathSandbox),
		ClientBasePublicURL(serverURL+apiPublicBasePathProd),
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

func TestClientAuth(t *testing.T) {
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
			w.Header().Add("Content-Type", mediaTypeApplicationJSON)

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

	client, clientMux, serverURL, clientTeardown, err := setupTestClient(token)
	if clientTeardown != nil {
		defer clientTeardown()
	}
	if !assert.NoError(err) {
		return
	}
	_, _ = client, clientMux

	assert.Equal(serverURL+apiBasePathSandbox, client.GetApiBaseURL())
	assert.Equal(serverURL+apiPublicBasePathProd, client.GetApiPublicBaseURL())
	assert.Equal(int64(1), authSeq, "must have generated 1 token")

	// All requests should have the Authorization header properly set
	clientMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		val := fmt.Sprintf("Bearer %s", buildAccessToken(authSeq-1))
		testHeader(t, r, "Authorization", val)
	})

	req1, err := client.newApiRequest(context.Background(), http.MethodGet, "test_auth1", nil, nil)
	if assert.NoError(err) {
		_, err := client.do(req1)
		assert.NoError(err)
	}

	assert.Equal(int64(1), authSeq, "a call with a valid token must not trigger token refresh")
	<-time.After(2 * time.Second)

	req2, err := client.newApiRequest(context.Background(), http.MethodGet, "test_auth2", nil, nil)
	if assert.NoError(err) {
		_, err := client.do(req2)
		assert.NoError(err)
	}

	assert.Equal(int64(2), authSeq, "after the token expired, it must be refreshed")
}
