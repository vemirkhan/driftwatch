package drift

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompare_NoDrift(t *testing.T) {
	d := NewDetector()
	expected := map[string]interface{}{
		"replicas": 3,
		"image":    "nginx:1.25",
	}
	actual := map[string]interface{}{
		"replicas": 3,
		"image":    "nginx:1.25",
	}

	report := d.Compare("web", "default", expected, actual)

	assert.False(t, report.Drifted)
	assert.Empty(t, report.Diffs)
	assert.Equal(t, "web", report.DeploymentName)
	assert.Equal(t, "default", report.Namespace)
}

func TestCompare_ValueChanged(t *testing.T) {
	d := NewDetector()
	expected := map[string]interface{}{"replicas": 3}
	actual := map[string]interface{}{"replicas": 1}

	report := d.Compare("web", "default", expected, actual)

	require.True(t, report.Drifted)
	require.Len(t, report.Diffs, 1)
	assert.Equal(t, "replicas", report.Diffs[0].Field)
	assert.Equal(t, 3, report.Diffs[0].Expected)
	assert.Equal(t, 1, report.Diffs[0].Actual)
}

func TestCompare_MissingField(t *testing.T) {
	d := NewDetector()
	expected := map[string]interface{}{"image": "nginx:1.25", "replicas": 2}
	actual := map[string]interface{}{"image": "nginx:1.25"}

	report := d.Compare("api", "prod", expected, actual)

	require.True(t, report.Drifted)
	require.Len(t, report.Diffs, 1)
	assert.Equal(t, "replicas", report.Diffs[0].Field)
	assert.Nil(t, report.Diffs[0].Actual)
}

func TestSummary_NoDrift(t *testing.T) {
	report := DriftReport{DeploymentName: "svc", Namespace: "staging", Drifted: false}
	assert.Contains(t, report.Summary(), "[OK]")
	assert.Contains(t, report.Summary(), "staging/svc")
}

func TestSummary_WithDrift(t *testing.T) {
	report := DriftReport{
		DeploymentName: "svc",
		Namespace:      "staging",
		Drifted:        true,
		Diffs: []FieldDiff{
			{Field: "replicas", Expected: 3, Actual: 1},
		},
	}
	summary := report.Summary()
	assert.Contains(t, summary, "[DRIFT]")
	assert.True(t, strings.Contains(summary, "replicas"))
}
