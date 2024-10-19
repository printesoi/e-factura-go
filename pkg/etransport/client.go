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

	xoauth2 "golang.org/x/oauth2"

	"github.com/printesoi/e-factura-go/pkg/client"
)

// ClientConfig is the config used to create a Client
type ClientConfig struct {
	// the client to use for making requests to the ANAF APIs protected with OAuth2.
	ApiClient *client.ApiClient
}

// ClientConfigOption allows gradually modifying a ClientConfig
type ClientConfigOption func(*ClientConfig)

// ClientApiClient sets the ApiClient to use for APIs protected with OAuth2.
func ClientApiClient(apiClient *client.ApiClient) ClientConfigOption {
	return func(c *ClientConfig) {
		c.ApiClient = apiClient
	}
}

// Client is a client that talks to ANAF e-transport APIs.
type Client struct {
	apiClient *client.ApiClient
}

// NewProductionClient creates a new basic Client for the ANAF e-transport production APIs.
func NewProductionClient(ctx context.Context, tokenSource xoauth2.TokenSource) (*Client, error) {
	apiClient, err := client.NewApiClient(
		client.ApiClientContext(ctx),
		client.ApiClientProductionEnvironment(true),
		client.ApiClientOAuth2TokenSource(tokenSource),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		apiClient: apiClient,
	}, nil
}

// NewSandboxClient creates a new basic Client for the ANAF e-transport sandbox(test) APIs.
func NewSandboxClient(ctx context.Context, tokenSource xoauth2.TokenSource) (*Client, error) {
	apiClient, err := client.NewApiClient(
		client.ApiClientContext(ctx),
		client.ApiClientSandboxEnvironment(true),
		client.ApiClientOAuth2TokenSource(tokenSource),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		apiClient: apiClient,
	}, nil
}

// NewClient allow for more control than NewProductionClient and NewSandboxClient
// by passing custom ApiClient to this Client.
func NewClient(opts ...ClientConfigOption) (*Client, error) {
	cfg := &ClientConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	return &Client{
		apiClient: cfg.ApiClient,
	}, nil
}
