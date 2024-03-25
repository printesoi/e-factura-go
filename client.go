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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

var (
	ErrInvalidClientOAuth2Config = errors.New("Invalid OAuth2Config provided")
	ErrInvalidClientOAuth2Token  = errors.New("Invalid Auth token provided")
)

const (
	Version = "v0.0.1-alpha"

	efacturaVersion  = "e-factura-go" + "/" + Version
	defaultUserAgent = efacturaVersion

	defaultApiANAF       = "https://api.anaf.ro"
	defaultApiPublicANAF = "https://webservicesp.anaf.ro"

	// apiBaseSandbox points to the sandbox (testing) version of the API
	apiBasePathSandbox = "/test/FCTEL/rest/"
	apiBaseSandbox     = defaultApiANAF + apiBasePathSandbox
	// apiBaseProd points to the production version of the API
	apiBasePathProd = "/prod/FCTEL/rest/"
	apiBaseProd     = defaultApiANAF + apiBasePathProd
	// apiPublicBaseProd points to the production version of the public
	// (unprotected) API.
	apiPublicBasePathProd = "/prod/FCTEL/rest/"
	apiPublicBaseProd     = defaultApiPublicANAF + apiPublicBasePathProd

	apiPathUpload                = "upload"
	apiPathMessageState          = "stareMesaj"
	apiPathMessageList           = "listaMesajeFactura"
	apiPathMessagePaginationList = "listaMesajePaginatieFactura"
	apiPathDownload              = "descarcare"

	webserviceAppPathValidate = "validare/%s"
	webserviceAppPathXMLToPDF = "transformare/%s"
)

// A Client manages communication with the ANAF e-factura APIs using OAuth2
// credentials.
type Client struct {
	apiBaseURL       *url.URL
	apiPublicBaseURL *url.URL
	userAgent        string

	oauth2Cfg    OAuth2Config
	initialToken *oauth2.Token

	apiClient *http.Client
}

// NewClient creates a new client using the provided config options.
func NewClient(ctx context.Context, opts ...ClientConfigOption) (*Client, error) {
	cfg := ClientConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	var apiBaseURL, apiPublicBaseURL string
	if cfg.BaseURL != nil {
		apiBaseURL = *cfg.BaseURL
	} else {
		apiBaseURL = getApiBase(cfg.Sandbox)
	}
	if cfg.BasePublicURL != nil {
		apiPublicBaseURL = *cfg.BasePublicURL
	} else {
		apiPublicBaseURL = apiPublicBaseProd
	}

	if !cfg.OAuth2Config.Valid() {
		return nil, ErrInvalidClientOAuth2Config
	}
	if !cfg.InitialToken.Valid() {
		return nil, ErrInvalidClientOAuth2Token
	}

	client := new(Client)
	client.userAgent = defaultUserAgent
	client.oauth2Cfg = cfg.OAuth2Config
	client.initialToken = cfg.InitialToken
	client.apiClient = cfg.OAuth2Config.Client(ctx, cfg.InitialToken)
	if cfg.UserAgent != nil {
		client.userAgent = *cfg.UserAgent
	}
	if _, err := client.withApiBaseURL(apiBaseURL, cfg.InsecureSkipVerify); err != nil {
		return client, err
	}
	if _, err := client.withApiPublicBaseURL(apiPublicBaseURL, cfg.InsecureSkipVerify); err != nil {
		return client, err
	}

	return client, nil
}

func (c *Client) withApiBaseURL(baseURL string, insecureSkipVerify bool) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return c, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		return c, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", baseURL)
	}

	c.apiBaseURL = base
	return c, nil
}

func (c *Client) withApiPublicBaseURL(baseURL string, insecureSkipVerify bool) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return c, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		return c, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", baseURL)
	}

	c.apiPublicBaseURL = base
	return c, nil
}

// GetApiBaseURL returns the base URL as string used by this client.
func (c *Client) GetApiBaseURL() string {
	return c.apiBaseURL.String()
}

// GetApiPublicBaseURL returns the base URL as string used by this client.
func (c *Client) GetApiPublicBaseURL() string {
	return c.apiPublicBaseURL.String()
}

func (c *Client) do(req *http.Request) (resp *http.Response, err error) {
	resp, err = c.apiClient.Do(req)
	if err == nil && !responseIsSuccess(resp.StatusCode) {
		err = newErrorResponse(resp, nil)
		return
	}
	return
}

func (c *Client) doApiUnmarshalXML(req *http.Request, response any) error {
	resp, err := c.do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if !responseBodyIsXML(resp.Header) {
		if responseBodyIsPlainText(resp.Header) {
			return newErrorResponseDetectType(resp)
		}
		return newErrorResponse(resp,
			fmt.Errorf("expected %s, got %s", mediaTypeApplicationXML, responseMediaType(resp.Header)))
	}
	if err := xmlUnmarshalReader(resp.Body, response); err != nil {
		return newErrorResponseParse(resp, err, false)
	}
	return nil
}

func (c *Client) doApiUnmarshalJSON(req *http.Request, destResponse any, cb func(*http.Response, any) error) error {
	resp, err := c.do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if !responseBodyIsJSON(resp.Header) {
		if responseBodyIsPlainText(resp.Header) {
			return newErrorResponseDetectType(resp)
		}
		return newErrorResponse(resp,
			fmt.Errorf("expected %s, got %s", mediaTypeApplicationJSON, responseMediaType(resp.Header)))
	}
	if err := jsonUnmarshalReader(resp.Body, destResponse); err != nil {
		return newErrorResponseParse(resp, err, false)
	}
	if cb != nil {
		return cb(resp, destResponse)
	}
	return nil
}

// RequestOption represents an option that can modify an http.Request.
type RequestOption func(req *http.Request)

// newRequest creates an API request. refURL is resolved relative to the given
// baseURL. The base URL should always has a trailing slash, while the relative
// URL should always be specified without a preceding slash.
func (c *Client) newRequest(ctx context.Context, method string, baseURL *url.URL,
	refURL string, query url.Values, body io.Reader, opts ...RequestOption,
) (*http.Request, error) {
	if !strings.HasSuffix(baseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", baseURL)
	}

	urlStr, err := buildURL(baseURL, refURL, query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, body)
	if err != nil {
		return req, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

// newApiRequest create an API request for the protected ANAF API.
func (c *Client) newApiRequest(ctx context.Context, method, refURL string,
	query url.Values, body io.Reader, opts ...RequestOption,
) (*http.Request, error) {
	return c.newRequest(ctx, method, c.apiBaseURL, refURL, query, body, opts...)
}

// newApiPublicRequest create an API request for the public ANAF API.
func (c *Client) newApiPublicRequest(ctx context.Context, method, refURL string,
	query url.Values, body io.Reader, opts ...RequestOption,
) (*http.Request, error) {
	return c.newRequest(ctx, method, c.apiPublicBaseURL, refURL, query, body, opts...)
}

func getApiBase(sandbox bool) string {
	if sandbox {
		return apiBaseSandbox
	}
	return apiBaseProd
}
