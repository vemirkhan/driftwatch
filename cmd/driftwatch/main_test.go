package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildBinary compiles the binary into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "driftwatch")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return binPath
}

func TestMain_MissingFlags(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(), "KUBECONFIG=/dev/null")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit when required flags missing")
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
		}
	}
}

func TestMain_InvalidChartPath(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin,
		"--deployment", "my-app",
		"--chart", "/nonexistent/values.yaml",
		"--namespace", "default",
	)
	cmd.Env = append(os.Environ(), "KUBECONFIG=/dev/null")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for invalid chart path")
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 0 {
			t.Error("expected non-zero exit code")
		}
	}
}

func TestMain_HelpFlag(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "--help")
	out, err := cmd.CombinedOutput()
	// --help exits with code 2 via flag package
	if err == nil {
		t.Fatal("expected non-zero exit for --help")
	}
	outStr := string(out)
	for _, keyword := range []string{"deployment", "chart", "output", "namespace"} {
		if !containsStr(outStr, keyword) {
			t.Errorf("expected help output to contain %q", keyword)
		}
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
