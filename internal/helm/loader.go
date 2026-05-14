package helm

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

// ChartValues holds the resolved default values from a Helm chart.
type ChartValues struct {
	ChartName string
	Version   string
	Values    map[string]interface{}
}

// Loader is responsible for loading Helm charts from the filesystem.
type Loader struct{}

// NewLoader creates a new Loader instance.
func NewLoader() *Loader {
	return &Loader{}
}

// LoadFromPath loads a Helm chart from the given directory or archive path
// and returns its resolved default values.
func (l *Loader) LoadFromPath(chartPath string) (*ChartValues, error) {
	if _, err := os.Stat(chartPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("chart path does not exist: %s", chartPath)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart from %s: %w", chartPath, err)
	}

	// CoalesceValues merges chart defaults with any provided overrides.
	// Here we pass nil to get pure chart defaults.
	coalesced, err := chartutil.CoalesceValues(chart, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to coalesce chart values: %w", err)
	}

	return &ChartValues{
		ChartName: chart.Metadata.Name,
		Version:   chart.Metadata.Version,
		Values:    coalesced.AsMap(),
	}, nil
}

// LoadFromValues accepts a pre-parsed values map (e.g., from a release secret)
// and wraps it in a ChartValues struct for uniform downstream handling.
func (l *Loader) LoadFromValues(name, version string, values map[string]interface{}) *ChartValues {
	if values == nil {
		values = make(map[string]interface{})
	}
	return &ChartValues{
		ChartName: name,
		Version:   version,
		Values:    values,
	}
}
