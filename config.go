package efactura

import (
	"golang.org/x/oauth2"
)

// ClientConfig is the config used to create a Client
type ClientConfig struct {
	// OAuth2Config is the OAuth2 config used for creating the http.Client that
	// autorefreshes the Token.
	OAuth2Config OAuth2Config
	// Token is the starting oauth2 Token (including the refresh token).
	// Until this library will support authentication with the SPV certificate,
	// this must always be provided.
	InitialToken *oauth2.Token
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

// ClientOAuth2Config sets the OAuth2 config
func ClientOAuth2Config(oauth2Cfg OAuth2Config) ClientConfigOption {
	return func(c *ClientConfig) {
		c.OAuth2Config = oauth2Cfg
	}
}

// ClientOAuth2InitialToken sets the initial OAuth2 Token
func ClientOAuth2InitialToken(token *oauth2.Token) ClientConfigOption {
	return func(c *ClientConfig) {
		c.InitialToken = token
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
