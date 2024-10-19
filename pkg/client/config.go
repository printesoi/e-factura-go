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
	"net/http"

	xoauth2 "golang.org/x/oauth2"

	"github.com/printesoi/e-factura-go/internal/ptr"
)

// baseClientConfig is the config used to create a baseClient
type baseClientConfig struct {
	// Base URL of the client.
	BaseURL string
	// http.Client to use for making requests.
	HttpClient *http.Client
	// User agent used when communicating with the ANAF API.
	UserAgent *string
	// Whether to skip the verification of the SSL certificate (default false).
	// Since this is a security risk, it should only be use with a custom
	// BaseURL in development/testing environments.
	InsecureSkipVerify bool
}

// baseClientConfigOption allows gradually modifying a baseClientConfig
type baseClientConfigOption func(*baseClientConfig)

// baseClientBaseURL sets the BaseURL to the given url.
func baseClientBaseURL(baseURL string) baseClientConfigOption {
	return func(c *baseClientConfig) {
		c.BaseURL = baseURL
	}
}

// baseClientHttpClient set the http.Client to use for making requests.
func baseClientHttpClient(client *http.Client) baseClientConfigOption {
	return func(c *baseClientConfig) {
		c.HttpClient = client
	}
}

// baseClientUserAgent sets the user agent used to communicate with the ANAF API.
func baseClientUserAgent(userAgent string) baseClientConfigOption {
	return func(c *baseClientConfig) {
		c.UserAgent = ptr.String(userAgent)
	}
}

// baseClientInsecureSkipVerify allows only setting InsecureSkipVerify. Please
// check the documentation for the InsecureSkipVerify field for a warning.
func baseClientInsecureSkipVerify(skipVerify bool) baseClientConfigOption {
	return func(c *baseClientConfig) {
		c.InsecureSkipVerify = skipVerify
	}
}

// PublicApiClientConfig is the config used to create a PublicApiClient
type PublicApiClientConfig struct {
	// Base URL of the ANAF public APIs. It is only useful in
	// development/testing environments.
	BaseURL *string
	// User agent used when communicating with the ANAF API.
	UserAgent *string
	// Whether to skip the verification of the SSL certificate (default false).
	// Since this is a security risk, it should only be use with a custom
	// BaseURL in development/testing environments.
	InsecureSkipVerify bool
}

// PublicApiClientConfigOption allows gradually modifying a PublicApiClientConfig
type PublicApiClientConfigOption func(*PublicApiClientConfig)

// PublicApiClientBaseURL sets the BaseURL to the given url. This should only
// be used when testing or using a custom endpoint for debugging.
func PublicApiClientBaseURL(baseURL string) PublicApiClientConfigOption {
	return func(c *PublicApiClientConfig) {
		c.BaseURL = ptr.String(baseURL)
	}
}

// PublicApiClientUserAgent sets the user agent used to communicate with the ANAF API.
func PublicApiClientUserAgent(userAgent string) PublicApiClientConfigOption {
	return func(c *PublicApiClientConfig) {
		c.UserAgent = ptr.String(userAgent)
	}
}

// PublicApiClientInsecureSkipVerify allows only setting InsecureSkipVerify. Please
// check the documentation for the InsecureSkipVerify field for a warning.
func PublicApiClientInsecureSkipVerify(skipVerify bool) PublicApiClientConfigOption {
	return func(c *PublicApiClientConfig) {
		c.InsecureSkipVerify = skipVerify
	}
}

// ApiClientConfig is the config used to create an ApiClient
type ApiClientConfig struct {
	// TokenSource is the token source used for generating OAuth2 tokens.
	// Until this library will support authentication with the SPV certificate,
	// this must always be provided.
	TokenSource xoauth2.TokenSource
	// Unless BaseURL is set, Sandbox controlls whether to use production
	// endpoints (if set to false) or test endpoints (if set to true).
	Sandbox bool
	// Context to use for creating the HTTP client. If not set,
	// context.Background will be used.
	Ctx context.Context
	// User agent used when communicating with the ANAF API.
	UserAgent *string
	// Base URL of the ANAF protected APIs. It is only useful in
	// development/testing environments.
	BaseURL *string
	// Whether to skip the verification of the SSL certificate (default false).
	// Since this is a security risk, it should only be use with a custom
	// BaseURL in development/testing environments.
	InsecureSkipVerify bool
}

// ApiClientConfigOption allows gradually modifying a ApiClientConfig
type ApiClientConfigOption func(*ApiClientConfig)

// ApiClientOAuth2TokenSource sets the token source to use for authorizing
// requests.
func ApiClientOAuth2TokenSource(tokenSource xoauth2.TokenSource) ApiClientConfigOption {
	return func(c *ApiClientConfig) {
		c.TokenSource = tokenSource
	}
}

// ApiClientSandboxEnvironment is the inverse of ApiClientProductionEnvironment:
// if called with sandbox=true sets the BaseURL to the sandbox URL,
// if called with sandbox=false sets the BaseURL to the production URL.
func ApiClientSandboxEnvironment(sandbox bool) ApiClientConfigOption {
	return func(c *ApiClientConfig) {
		c.Sandbox = sandbox
	}
}

// ApiClientProductionEnvironment is the inverse of ApiClientSandboxEnvironment
// if called with prod=true sets the BaseURL to the production URL,
// if called with prod=false sets the BaseURL to the sandbox URL.
func ApiClientProductionEnvironment(prod bool) ApiClientConfigOption {
	return func(c *ApiClientConfig) {
		c.Sandbox = !prod
	}
}

// ApiClientContext sets the Context to use for building the http.Client.
func ApiClientContext(ctx context.Context) ApiClientConfigOption {
	return func(c *ApiClientConfig) {
		c.Ctx = ctx
	}
}

// ApiClientUserAgent sets the user agent used to communicate with the ANAF API.
func ApiClientUserAgent(userAgent string) ApiClientConfigOption {
	return func(c *ApiClientConfig) {
		c.UserAgent = ptr.String(userAgent)
	}
}

// ApiClientBaseURL sets the BaseURL to the given url. This should only
// be used when testing or using a custom endpoint for debugging/testing.
func ApiClientBaseURL(baseURL string) ApiClientConfigOption {
	return func(c *ApiClientConfig) {
		c.BaseURL = ptr.String(baseURL)
	}
}

// ApiClientInsecureSkipVerify allows only setting InsecureSkipVerify. Please
// check the documentation for the InsecureSkipVerify field for a warning.
func ApiClientInsecureSkipVerify(skipVerify bool) ApiClientConfigOption {
	return func(c *ApiClientConfig) {
		c.InsecureSkipVerify = skipVerify
	}
}
