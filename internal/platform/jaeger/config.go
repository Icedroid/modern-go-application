package jaeger

import "github.com/pkg/errors"

// Config holds information necessary for sending trace to Jaeger.
type Config struct {
	// Endpoint is the Jaeger HTTP Thrift endpoint.
	// For example, http://localhost:14268.
	Endpoint string

	// AgentEndpoint instructs exporter to send spans to Jaeger agent at this address.
	// For example, localhost:6831.
	AgentEndpoint string

	// Username to be used if basic auth is required.
	// Optional.
	Username string

	// Password to be used if basic auth is required.
	// Optional.
	Password string
}

// Validate checks that the configuration is valid.
func (c Config) Validate() error {
	if c.Endpoint == "" && c.AgentEndpoint == "" {
		return errors.New("either endpoint or agent endpoint must be configured")
	}

	return nil
}
