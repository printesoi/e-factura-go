package efactura

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

var (
	ErrInvalidClientOAuth2Config = errors.New("Invalid OAuth2Config provided")
	ErrInvalidClientOAuth2Token  = errors.New("Invalid Auth token provided")
)

const (
	// apiBaseSandbox points to the sandbox (testing) version of the API
	apiBaseSandbox           = "https://api.anaf.ro/test/FCTEL/rest/"
	webserviceAppBaseSandbox = "https://webserviceapl.anaf.ro/test/FCTEL/rest/"
	// apiBaseProd points to the production version of the API
	apiBaseProd           = "https://api.anaf.ro/prod/FCTEL/rest/"
	webserviceAppBaseProd = "https://webserviceapl.anaf.ro/prod/FCTEL/rest/"
	webserviceSpBaseProd  = "https://webservicesp.anaf.ro/prod/FCTEL/rest/"

	apiPathUpload                = "/upload"
	apiPathMessageState          = "/stareMesaj"
	apiPathMessageList           = "/listaMesajeFactura"
	apiPathMessagePaginationList = "/listaMesajePaginatieFactura"
	apiPathDownload              = "/descarcare"

	webserviceAppPathValidate = "/validare/%s"
	webserviceAppPathXmlToPdf = "/transformare/%s"
)

// Client is a client for the ANAF e-factura APIs using OAuth2 credentials.
type Client struct {
	sandbox    bool
	apiBaseUrl *url.URL

	oauth2Cfg OAuth2Config
	token     oauth2.Token

	apiClient *http.Client
}

// NewClient creates a new client using the provided config options.
func NewClient(ctx context.Context, opts ...ClientConfigOption) (*Client, error) {
	cfg := ClientConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	var baseUrl string
	if cfg.BaseUrl != nil {
		baseUrl = *cfg.BaseUrl
	} else {
		baseUrl = getApiBase(cfg.Sandbox)
	}

	if !cfg.OAuth2Config.Valid() {
		return nil, ErrInvalidClientOAuth2Config
	}
	if cfg.Token == nil || !cfg.Token.Valid() {
		return nil, ErrInvalidClientOAuth2Token
	}

	return (&Client{
		sandbox:   cfg.Sandbox,
		apiClient: cfg.OAuth2Config.Client(ctx, cfg.Token),
	}).withApiBaseURL(baseUrl, cfg.InsecureSkipVerify)
}

// GetApiBaseUrl returns the base URL as string used by this client.
func (c *Client) GetApiBaseUrl() string {
	return c.apiBaseUrl.String()
}

func (c *Client) withApiBaseURL(baseURL string, insecureSkipVerify bool) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return c, err
	}

	c.apiBaseUrl = base
	return c, nil
}

func (c *Client) buildApiUrl(path string, query url.Values) string {
	return buildUrl(c.apiBaseUrl, path, query)
}

func (c *Client) doApi(req *http.Request) (body []byte, statusCode int, headers http.Header, err error) {
	var resp *http.Response
	resp, err = c.apiClient.Do(req)
	if resp != nil {
		statusCode = resp.StatusCode
		headers = resp.Header

		if resp.Body != nil {
			defer resp.Body.Close()
			body, err = io.ReadAll(resp.Body)
		}

		if err == nil && !responseIsSuccess(resp.StatusCode) {
			err = newErrorResponse(resp, nil)
			return
		}
	}
	return
}

func (c *Client) doApiUnmarshalXML(req *http.Request, response any) error {
	resp, err := c.apiClient.Do(req)
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

func getApiBase(sandbox bool) string {
	if sandbox {
		return apiBaseSandbox
	}
	return apiBaseProd
}

func buildUrl(base *url.URL, path string, query url.Values) string {
	u := base
	u.Path, _ = url.JoinPath(base.Path, path)
	u.RawQuery = query.Encode()
	return u.String()
}

func buildParseUrl(base string, path string, query url.Values) (string, error) {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	retUrl := buildUrl(baseUrl, path, query)
	return retUrl, nil
}

func newRequest(ctx context.Context, method, url string, payload []byte) (*http.Request, error) {
	var buf io.Reader
	if payload != nil {
		buf = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return req, err
	}

	return req, nil
}

// This is a copy of the drainBody from src/net/http/httputil/dump.go
func drainBody(b io.ReadCloser) (body []byte, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		return nil, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return buf.Bytes(), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func peekResponseBody(r *http.Response) (body []byte, err error) {
	body, r.Body, err = drainBody(r.Body)
	return
}

func peekRequestBody(r *http.Request) (body []byte, err error) {
	body, r.Body, err = drainBody(r.Body)
	return
}

func responseMediaType(headers http.Header) (mediaType string) {
	mediaType, _, _ = mime.ParseMediaType(headers.Get("Content-Type"))
	return
}

func responseBodyIsJSON(headers http.Header) bool {
	return responseMediaType(headers) == "application/json"
}

func responseBodyIsXML(headers http.Header) bool {
	switch responseMediaType(headers) {
	case "application/xml", "text/xml":
		return true
	}
	return false
}

func responseIsSuccess(status int) bool {
	return status >= 200 && status < 300
}

func jsonUnmarshalReader(r io.Reader, v any) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func xmlUnmarshalReader(r io.Reader, v any) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}
