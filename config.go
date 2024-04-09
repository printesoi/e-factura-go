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
	xoauth2 "golang.org/x/oauth2"
)

// ClientConfig is the config used to create a Client
type ClientConfig struct {
	// TokenSource is the token source used for generating OAuth2 tokens.
	// Until this library will support authentication with the SPV certificate,
	// this must always be provided.
	TokenSource xoauth2.TokenSource
	// Unless BaseURL is set, Sandbox controlls whether to use production
	// endpoints (if set to false) or test endpoints (if set to true).
	Sandbox bool
	// User agent used when communicating with the ANAF API.
	UserAgent *string
	// Base URL of the ANAF e-factura protected APIs. It is only useful in
	// development/testing environments.
	BaseURL *string
	// Base URL of the ANAF e-factura public(unprotected) APIs. It is only
	// useful in development/testing environments.
	BasePublicURL *string
	// Whether to skip the verification of the SSL certificate (default false).
	// Since this is a security risk, it should only be use with a custom
	// BaseURL in development/testing environments.
	InsecureSkipVerify bool
}

// ClientConfigOption allows gradually modifying a ClientConfig
type ClientConfigOption func(*ClientConfig)

// ClientOAuth2TokenSource sets the token source to use.
func ClientOAuth2TokenSource(tokenSource xoauth2.TokenSource) ClientConfigOption {
	return func(c *ClientConfig) {
		c.TokenSource = tokenSource
	}
}

// ClientSandboxEnvironment(true) set the BaseURL to the sandbox URL
func ClientSandboxEnvironment(sandbox bool) ClientConfigOption {
	return func(c *ClientConfig) {
		c.Sandbox = sandbox
	}
}

// ClientProductionEnvironment(true) set the BaseURL to the production URL
func ClientProductionEnvironment(prod bool) ClientConfigOption {
	return func(c *ClientConfig) {
		c.Sandbox = !prod
	}
}

// ClientBaseURL sets the BaseURL to the given url. This should only be used
// when testing or using a custom endpoint for debugging.
func ClientBaseURL(baseURL string) ClientConfigOption {
	return func(c *ClientConfig) {
		c.BaseURL = ptrfyString(baseURL)
	}
}

// ClientBasePublicURL sets the BaseURL to the given url. This should only be
// used when testing or using a custom endpoint for debugging.
func ClientBasePublicURL(baseURL string) ClientConfigOption {
	return func(c *ClientConfig) {
		c.BasePublicURL = ptrfyString(baseURL)
	}
}

// ClientUserAgent sets the user agent used to communicate with the ANAF API.
func ClientUserAgent(userAgent string) ClientConfigOption {
	return func(c *ClientConfig) {
		c.UserAgent = ptrfyString(userAgent)
	}
}

// ClientInsecureSkipVerify allows only setting InsecureSkipVerify. Please
// check the documentation for the InsecureSkipVerify field for a warning.
func ClientInsecureSkipVerify(skipVerify bool) ClientConfigOption {
	return func(c *ClientConfig) {
		c.InsecureSkipVerify = skipVerify
	}
}
