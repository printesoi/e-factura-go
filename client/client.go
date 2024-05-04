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
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	xoauth2 "golang.org/x/oauth2"

	"github.com/printesoi/e-factura-go/constants"
	ierrors "github.com/printesoi/e-factura-go/internal/errors"
	api_helpers "github.com/printesoi/e-factura-go/internal/helpers/api"
	"github.com/printesoi/e-factura-go/xml"
)

const (
	defaultUserAgent = "anaf-api-go" + "/" + constants.Version
)

// baseClient is a HTTP client for the ANAF APIs. It's embedded in a ApiClient
// or PublicApiClient.
type baseClient struct {
	baseURL    *url.URL
	userAgent  string
	httpClient *http.Client
	wg         sync.WaitGroup
}

// newBaseClient creates a new baseClient using the provided config options.
func newBaseClient(opts ...baseClientConfigOption) (*baseClient, error) {
	cfg := baseClientConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	client := new(baseClient)
	client.userAgent = defaultUserAgent
	if cfg.UserAgent != nil {
		client.userAgent = *cfg.UserAgent
	}
	if cfg.HttpClient != nil {
		client.httpClient = cfg.HttpClient
	} else {
		client.httpClient = &http.Client{}
	}

	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return client, err
	}
	if !strings.HasSuffix(baseURL.Path, "/") {
		return client, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", cfg.BaseURL)
	}
	client.baseURL = baseURL
	return client, nil
}

// Do sends the given HTTP request and returns an HTTP response. A non-200
// response results in an *errors.ErrorResponse error.
func (c *baseClient) Do(req *http.Request) (resp *http.Response, err error) {
	c.wg.Add(1)
	defer c.wg.Done()

	resp, err = c.httpClient.Do(req)
	if err == nil && !api_helpers.ResponseIsSuccess(resp.StatusCode) {
		err = ierrors.NewErrorResponse(resp, nil)
		return
	}
	return
}

// DoUnmarshalXML sends the given HTTP request and expects an XML response
// which is unmarshalled into response. A non-200 response results in an
// *errors.ErrorResponse error. If response body if not application/xml, we try
// to autodetect if request limits are exceeded.
func (c *baseClient) DoUnmarshalXML(req *http.Request, response any) error {
	resp, err := c.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if !api_helpers.ResponseBodyIsXML(resp.Header) {
		if api_helpers.ResponseBodyIsPlainText(resp.Header) {
			return ierrors.NewErrorResponseDetectType(resp)
		}
		return ierrors.NewErrorResponse(resp,
			fmt.Errorf("expected %s, got %s", api_helpers.MediaTypeApplicationXML, api_helpers.ResponseMediaType(resp.Header)))
	}
	if err := xml.UnmarshalReaderXML(resp.Body, response); err != nil {
		return ierrors.NewErrorResponseParse(resp, err, false)
	}
	return nil
}

// DoUnmarshalJSON sends the given HTTP request and expects a JSON response
// which is unmarshalled into response. A non-200 response results in an
// *errors.ErrorResponse error. If response body if not application/json, we try
// to autodetect if request limits are exceeded.
func (c *baseClient) DoUnmarshalJSON(req *http.Request, destResponse any, cb func(*http.Response, any) error) error {
	resp, err := c.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if !api_helpers.ResponseBodyIsJSON(resp.Header) {
		if api_helpers.ResponseBodyIsPlainText(resp.Header) {
			return ierrors.NewErrorResponseDetectType(resp)
		}
		return ierrors.NewErrorResponse(resp,
			fmt.Errorf("expected %s, got %s", api_helpers.MediaTypeApplicationJSON, api_helpers.ResponseMediaType(resp.Header)))
	}
	if err := api_helpers.UnmarshalReaderJSON(resp.Body, destResponse); err != nil {
		return ierrors.NewErrorResponseParse(resp, err, false)
	}
	if cb != nil {
		return cb(resp, destResponse)
	}
	return nil
}

// Wait wait for all requests for finish
func (c *baseClient) Wait() {
	c.wg.Wait()
}

// RequestOption represents an option that can modify an http.Request.
type RequestOption func(req *http.Request)

// NewRequest creates an API request. refURL is resolved relative to the client
// baseURL. The relative URL should always be specified without a preceding slash.
func (c *baseClient) NewRequest(ctx context.Context, method string,
	refURL string, query url.Values, body io.Reader, opts ...RequestOption,
) (*http.Request, error) {
	urlStr, err := api_helpers.BuildURL(c.baseURL, refURL, query)
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

// PublicApiClient is a client that interacts with ANAF public APIs (no OAuth2
// authentication is performed automatically by this client).
type PublicApiClient struct {
	*baseClient
}

// NewPublicApiClient creates a new PublicApiClient using the provided config options.
func NewPublicApiClient(opts ...PublicApiClientConfigOption) (*PublicApiClient, error) {
	cfg := PublicApiClientConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	baseURL := constants.PublicApiBaseURL
	if cfg.BaseURL != nil {
		baseURL = *cfg.BaseURL
	}
	baseOpts := []baseClientConfigOption{
		baseClientBaseURL(baseURL),
	}
	if cfg.UserAgent != nil {
		baseOpts = append(baseOpts, baseClientUserAgent(*cfg.UserAgent))
	}
	if cfg.InsecureSkipVerify {
		baseOpts = append(baseOpts, baseClientInsecureSkipVerify(cfg.InsecureSkipVerify))
	}
	baseClient, err := newBaseClient(baseOpts...)
	if err != nil {
		return nil, err
	}

	return &PublicApiClient{
		baseClient: baseClient,
	}, nil
}

// ApiClient is a client that interacts with ANAF protected APIs (APIs that
// required OAuth2 authentication).
type ApiClient struct {
	*baseClient
}

// NewApiClient creates a new ApiClient using the provided config options.
func NewApiClient(opts ...ApiClientConfigOption) (*ApiClient, error) {
	cfg := ApiClientConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.TokenSource == nil {
		return nil, errors.New("invalid token source for client")
	}

	var baseURL string
	if cfg.BaseURL != nil {
		baseURL = *cfg.BaseURL
	} else {
		baseURL = getApiBase(cfg.Sandbox)
	}

	ctx := cfg.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	baseOpts := []baseClientConfigOption{
		baseClientBaseURL(baseURL),
		baseClientHttpClient(xoauth2.NewClient(ctx, cfg.TokenSource)),
		baseClientInsecureSkipVerify(cfg.InsecureSkipVerify),
	}
	if cfg.UserAgent != nil {
		baseOpts = append(baseOpts, baseClientUserAgent(*cfg.UserAgent))
	}
	if cfg.InsecureSkipVerify {
		baseOpts = append(baseOpts, baseClientInsecureSkipVerify(cfg.InsecureSkipVerify))
	}
	baseClient, err := newBaseClient(baseOpts...)
	if err != nil {
		return nil, err
	}

	return &ApiClient{
		baseClient: baseClient,
	}, nil
}

func getApiBase(sandbox bool) string {
	if sandbox {
		return constants.ApiBaseSandbox
	}
	return constants.ApiBaseProd
}
