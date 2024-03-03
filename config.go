package efactura

// ClientConfig is the config used to create a Client
type ClientConfig struct {
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

// ClientConfigOption allows gradually modifying the config
type ClientConfigOption func(*ClientConfig)

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
