package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/report"
)

func TestWriteText_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(&buf, report.FormatText)
	if err := f.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteText_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(&buf, report.FormatText)
	entries := []drift.DiffEntry{
		{Key: "image.tag", Type: "changed", ChartValue: "1.0.0", LiveValue: "1.1.0"},
		{Key: "replicaCount", Type: "missing", ChartValue: 3, LiveValue: nil},
	}
	if err := f.Write(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Drift detected (2") {
		t.Errorf("expected drift header, got: %s", out)
	}
	if !strings.Contains(out, "image.tag") {
		t.Errorf("expected image.tag in output, got: %s", out)
	}
	if !strings.Contains(out, "replicaCount") {
		t.Errorf("expected replicaCount in output, got: %s", out)
	}
}

func TestWriteJSON_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(&buf, report.FormatJSON)
	if err := f.Write([]drift.DiffEntry{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"drift":false`) {
		t.Errorf("expected drift:false in JSON, got: %s", out)
	}
}

func TestWriteJSON_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(&buf, report.FormatJSON)
	entries := []drift.DiffEntry{
		{Key: "image.repository", Type: "changed", ChartValue: "nginx", LiveValue: "apache"},
	}
	if err := f.Write(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"drift":true`) {
		t.Errorf("expected drift:true in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"image.repository"`) {
		t.Errorf("expected key in JSON output, got: %s", out)
	}
}
