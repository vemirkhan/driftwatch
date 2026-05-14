package helm

import (
	"testing"
)

func TestLoadFromValues_Basic(t *testing.T) {
	loader := NewLoader()

	values := map[string]interface{}{
		"replicaCount": 3,
		"image": map[string]interface{}{
			"repository": "nginx",
			"tag":        "1.25",
		},
	}

	cv := loader.LoadFromValues("my-chart", "1.0.0", values)

	if cv.ChartName != "my-chart" {
		t.Errorf("expected ChartName 'my-chart', got '%s'", cv.ChartName)
	}
	if cv.Version != "1.0.0" {
		t.Errorf("expected Version '1.0.0', got '%s'", cv.Version)
	}
	if cv.Values["replicaCount"] != 3 {
		t.Errorf("expected replicaCount 3, got %v", cv.Values["replicaCount"])
	}
}

func TestLoadFromValues_NilValues(t *testing.T) {
	loader := NewLoader()

	cv := loader.LoadFromValues("empty-chart", "0.1.0", nil)

	if cv.Values == nil {
		t.Error("expected non-nil Values map when nil is passed")
	}
	if len(cv.Values) != 0 {
		t.Errorf("expected empty Values map, got %d entries", len(cv.Values))
	}
}

func TestLoadFromValues_NestedValues(t *testing.T) {
	loader := NewLoader()

	values := map[string]interface{}{
		"service": map[string]interface{}{
			"type": "ClusterIP",
			"port": 80,
		},
	}

	cv := loader.LoadFromValues("svc-chart", "2.3.1", values)

	svc, ok := cv.Values["service"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'service' to be a map")
	}
	if svc["type"] != "ClusterIP" {
		t.Errorf("expected service.type 'ClusterIP', got '%v'", svc["type"])
	}
}

func TestLoadFromPath_NonExistent(t *testing.T) {
	loader := NewLoader()

	_, err := loader.LoadFromPath("/tmp/does-not-exist-chart-xyz")
	if err == nil {
		t.Error("expected error for non-existent chart path, got nil")
	}
}
