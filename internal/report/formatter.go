// Package report provides formatting and output utilities for drift detection results.
package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/driftwatch/internal/drift"
)

// Format represents the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes drift comparison results to an output writer.
type Formatter struct {
	w      io.Writer
	format Format
}

// NewFormatter creates a new Formatter writing to w in the given format.
func NewFormatter(w io.Writer, format Format) *Formatter {
	return &Formatter{w: w, format: format}
}

// Write outputs the drift results to the underlying writer.
func (f *Formatter) Write(results []drift.DiffEntry) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(results)
	default:
		return f.writeText(results)
	}
}

func (f *Formatter) writeText(results []drift.DiffEntry) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(f.w, "✓ No drift detected.")
		return err
	}

	_, err := fmt.Fprintf(f.w, "⚠ Drift detected (%d difference(s)):\n", len(results))
	if err != nil {
		return err
	}

	for _, entry := range results {
		line := formatTextEntry(entry)
		if _, err := fmt.Fprintln(f.w, line); err != nil {
			return err
		}
	}
	return nil
}

func formatTextEntry(e drift.DiffEntry) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  [%s] %s", strings.ToUpper(e.Type), e.Key))
	switch e.Type {
	case "changed":
		sb.WriteString(fmt.Sprintf(": chart=%v, live=%v", e.ChartValue, e.LiveValue))
	case "missing":
		sb.WriteString(fmt.Sprintf(": expected=%v (not found in live)", e.ChartValue))
	case "extra":
		sb.WriteString(fmt.Sprintf(": live=%v (not in chart)", e.LiveValue))
	}
	return sb.String()
}

func (f *Formatter) writeJSON(results []drift.DiffEntry) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(f.w, `{"drift":false,"differences":[]}`)
		return err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`{"drift":true,"differences":[`))
	for i, e := range results {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(
			`{"type":%q,"key":%q,"chart_value":%q,"live_value":%q}`,
			e.Type, e.Key, fmt.Sprintf("%v", e.ChartValue), fmt.Sprintf("%v", e.LiveValue),
		))
	}
	sb.WriteString(`]}`)
	_, err := fmt.Fprintln(f.w, sb.String())
	return err
}
