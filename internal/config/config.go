// Package config provides configuration loading and validation for driftwatch.
package config

import (
	"errors"
	"fmt"
	"os"
)

// Config holds the runtime configuration for a driftwatch scan.
type Config struct {
	// Namespace to scan; empty string means all namespaces.
	Namespace string
	// DeploymentName filters to a single deployment when non-empty.
	DeploymentName string
	// ChartPath is the path to the Helm chart directory or values file.
	ChartPath string
	// OutputFormat is either "text" or "json".
	OutputFormat string
	// KubeContext is the kubectl context name to use; empty uses current context.
	KubeContext string
	// FailOnDrift causes a non-zero exit code when drift is detected.
	FailOnDrift bool
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		OutputFormat: "text",
		FailOnDrift:  false,
	}
}

// Validate checks that the Config contains the required fields and valid values.
func (c *Config) Validate() error {
	if c.ChartPath == "" {
		return errors.New("chart path must not be empty")
	}

	if _, err := os.Stat(c.ChartPath); err != nil {
		return fmt.Errorf("chart path %q is not accessible: %w", c.ChartPath, err)
	}

	switch c.OutputFormat {
	case "text", "json":
		// valid
	case "":
		c.OutputFormat = "text"
	default:
		return fmt.Errorf("unsupported output format %q: must be \"text\" or \"json\"", c.OutputFormat)
	}

	return nil
}
