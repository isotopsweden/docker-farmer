package handlers

import (
	"errors"

	"github.com/tsuru/config"
)

// getDomain will return the domain from the configuration or a error.
func getDomain() (string, error) {
	// Check domain.
	domain, err := config.GetString("domain")

	if err != nil {
		return "", err
	}

	if domain == "" {
		return "", errors.New("Domain environment variable is missing")
	}

	return domain, nil
}
