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
	Version = "v0.0.1"

	defaultUserAgent     = "e-factura-go" + "/" + Version
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
	apiPublicBaseProd     = defaultApiPublicANAF + "/prod/FCTEL/rest/"

	apiPathUpload                = "upload"
	apiPathMessageState          = "stareMesaj"
	apiPathMessageList           = "listaMesajeFactura"
	apiPathMessagePaginationList = "listaMesajePaginatieFactura"
	apiPathDownload              = "descarcare"

	webserviceAppPathValidate = "validare/%s"
	webserviceAppPathXmlToPdf = "transformare/%s"
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

	c.apiBaseURL = base
	return c, nil
}

func (c *Client) withApiPublicBaseURL(baseURL string, insecureSkipVerify bool) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return c, err
	}

	c.apiPublicBaseURL = base
	return c, nil
}

// GetApiBaseUrl returns the base URL as string used by this client.
func (c *Client) GetApiBaseUrl() string {
	return c.apiBaseURL.String()
}

// GetApiPublicBaseUrl returns the base URL as string used by this client.
func (c *Client) GetApiPublicBaseUrl() string {
	return c.apiPublicBaseURL.String()
}

func (c *Client) do(req *http.Request) (resp *http.Response, err error) {
	resp, err = c.apiClient.Do(req)
	c.debugRequest(req, resp)
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
		return newErrorResponse(resp,
			fmt.Errorf("expected application/xml, got %s", responseMediaType(resp.Header)))
	}
	return xmlUnmarshalReader(resp.Body, response)
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
