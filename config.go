package efactura

import (
	"golang.org/x/oauth2"
)

// ClientConfig is the config used to create a Client
type ClientConfig struct {
	// OAuth2Config is the OAuth2 config used for creating the http.Client that
	// autorefreshes the Token.
	OAuth2Config OAuth2Config
	// Token is the starting oauth2 Token, until this library will support
	// authentication with the SPV certificate, this must always be provided.
	Token *oauth2.Token
	// Unless BaseUrl is set, Sandbox controlls whether to use production
	// endpoints (if set to false) or test endpoints (if set to true).
	Sandbox bool
	// Base URL of the ANAF e-factura APIs. It is only useful in
	// development/testing environments.
	BaseUrl *string
	// Whether to skip the verification of the SSL certificate (default false).
	// Since this is a security risk, it should only be use with a custom
	// BaseUrl in development/testing environments.
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

// ClientOAuth2Token sets the initial OAuth2 Token
func ClientOAuth2Token(token *oauth2.Token) ClientConfigOption {
	return func(c *ClientConfig) {
		c.Token = token
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

// ClientBaseUrl sets the BaseURL to the given url. This should only be used
// when testing or using a custom endpoint for debugging.
func ClientBaseUrl(baseUrl string) ClientConfigOption {
	return func(c *ClientConfig) {
		c.BaseUrl = &baseUrl
	}
}

// ClientInsecureSkipVerify allows only setting InsecureSkipVerify. Please
// check the documentation for the InsecureSkipVerify field for a warning.
func ClientInsecureSkipVerify(skipVerify bool) ClientConfigOption {
	return func(c *ClientConfig) {
		c.InsecureSkipVerify = skipVerify
	}
}
