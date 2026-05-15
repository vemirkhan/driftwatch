package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.OutputFormat != "text" {
		t.Errorf("expected default OutputFormat \"text\", got %q", cfg.OutputFormat)
	}
	if cfg.FailOnDrift {
		t.Error("expected FailOnDrift to be false by default")
	}
}

func TestValidate_MissingChartPath(t *testing.T) {
	cfg := config.DefaultConfig()
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing chart path, got nil")
	}
}

func TestValidate_NonExistentChartPath(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.ChartPath = "/nonexistent/path/to/chart"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for non-existent chart path, got nil")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	dir := t.TempDir()
	valuesFile := filepath.Join(dir, "values.yaml")
	if err := os.WriteFile(valuesFile, []byte("key: value\n"), 0o644); err != nil {
		t.Fatalf("failed to create temp values file: %v", err)
	}

	cfg := config.DefaultConfig()
	cfg.ChartPath = valuesFile
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestValidate_InvalidOutputFormat(t *testing.T) {
	dir := t.TempDir()
	cfg := config.DefaultConfig()
	cfg.ChartPath = dir
	cfg.OutputFormat = "yaml"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid output format, got nil")
	}
}

func TestValidate_EmptyOutputFormatDefaultsToText(t *testing.T) {
	dir := t.TempDir()
	cfg := config.DefaultConfig()
	cfg.ChartPath = dir
	cfg.OutputFormat = ""
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("expected OutputFormat to default to \"text\", got %q", cfg.OutputFormat)
	}
}
