// Package config defines the Config struct used to configure a driftwatch
// scan at runtime.
//
// A Config specifies which Kubernetes namespace and deployment to inspect,
// the path to the Helm chart or values file to compare against, the desired
// output format, and behavioural flags such as whether to exit with a
// non-zero status code when drift is detected.
//
// Typical usage:
//
//	cfg := config.DefaultConfig()
//	cfg.ChartPath = "/path/to/chart"
//	cfg.Namespace = "production"
//	if err := cfg.Validate(); err != nil {
//	    log.Fatal(err)
//	}
package config
