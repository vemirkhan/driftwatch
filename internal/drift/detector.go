package drift

import (
	"fmt"
	"reflect"
)

// FieldDiff represents a single field that has drifted.
type FieldDiff struct {
	Field    string
	Expected interface{}
	Actual   interface{}
}

// DriftReport holds the result of comparing a live deployment against
// its Helm chart definition.
type DriftReport struct {
	DeploymentName string
	Namespace      string
	Drifted        bool
	Diffs          []FieldDiff
}

// Detector compares live Kubernetes resource values against expected
// values sourced from a rendered Helm chart.
type Detector struct{}

// NewDetector creates a new Detector instance.
func NewDetector() *Detector {
	return &Detector{}
}

// Compare takes a map of expected values (from Helm) and actual values
// (from the live cluster) and returns a DriftReport.
func (d *Detector) Compare(name, namespace string, expected, actual map[string]interface{}) DriftReport {
	report := DriftReport{
		DeploymentName: name,
		Namespace:      namespace,
	}

	for key, expectedVal := range expected {
		actualVal, exists := actual[key]
		if !exists {
			report.Diffs = append(report.Diffs, FieldDiff{
				Field:    key,
				Expected: expectedVal,
				Actual:   nil,
			})
			continue
		}
		if !reflect.DeepEqual(expectedVal, actualVal) {
			report.Diffs = append(report.Diffs, FieldDiff{
				Field:    key,
				Expected: expectedVal,
				Actual:   actualVal,
			})
		}
	}

	report.Drifted = len(report.Diffs) > 0
	return report
}

// Summary returns a human-readable summary of the drift report.
func (r DriftReport) Summary() string {
	if !r.Drifted {
		return fmt.Sprintf("[OK] %s/%s: no drift detected", r.Namespace, r.DeploymentName)
	}
	msg := fmt.Sprintf("[DRIFT] %s/%s: %d field(s) differ\n", r.Namespace, r.DeploymentName, len(r.Diffs))
	for _, diff := range r.Diffs {
		msg += fmt.Sprintf("  - %s: expected=%v actual=%v\n", diff.Field, diff.Expected, diff.Actual)
	}
	return msg
}
