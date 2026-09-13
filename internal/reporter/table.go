package reporter

import (
	"fmt"
	"strings"

	"github.com/Pixelcity-dev/Deepsec/internal/core"
)

type TableReporter struct{}

func (r *TableReporter) Name() string {
	return "table"
}

func (r *TableReporter) Extension() string {
	return ""
}

func (r *TableReporter) Generate(results []core.ScanResult, opts ReportOptions) ([]byte, error) {
	var sb strings.Builder

	totalFindings := 0
	for _, result := range results {
		totalFindings += len(result.Findings)
	}

	sb.WriteString("\n")
	sb.WriteString("╔══════════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║                    DeepSec Scan Results                      ║\n")
	sb.WriteString("╚══════════════════════════════════════════════════════════════╝\n\n")

	if totalFindings == 0 {
		sb.WriteString("✅ No security issues found!\n\n")
		return []byte(sb.String()), nil
	}

	sb.WriteString(fmt.Sprintf("Found %d security issues:\n\n", totalFindings))

	bySeverity := make(map[string]int)
	byType := make(map[string]int)

	for _, result := range results {
		for _, f := range result.Findings {
			bySeverity[f.Severity.String()]++
			byType[string(f.ScanType)]++
		}
	}

	sb.WriteString("Summary:\n")
	sb.WriteString("────────\n")
	for severity, count := range bySeverity {
		if opts.Color {
			switch severity {
			case "CRITICAL":
				sb.WriteString(fmt.Sprintf("  \033[1;31m%-12s\033[0m %d\n", severity, count))
			case "HIGH":
				sb.WriteString(fmt.Sprintf("  \033[31m%-12s\033[0m %d\n", severity, count))
			case "MEDIUM":
				sb.WriteString(fmt.Sprintf("  \033[93m%-12s\033[0m %d\n", severity, count))
			case "LOW":
				sb.WriteString(fmt.Sprintf("  \033[33m%-12s\033[0m %d\n", severity, count))
			default:
				sb.WriteString(fmt.Sprintf("  %-12s %d\n", severity, count))
			}
		} else {
			sb.WriteString(fmt.Sprintf("  %-12s %d\n", severity, count))
		}
	}
	sb.WriteString("\n")

	for scanType, count := range byType {
		sb.WriteString(fmt.Sprintf("  %-15s %d issues\n", scanType, count))
	}
	sb.WriteString("\n")

	for _, result := range results {
		if len(result.Findings) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("Scanner: %s\n", result.Scanner))
		sb.WriteString(strings.Repeat("─", 60) + "\n")

		for _, f := range result.Findings {
			severity := f.Severity.String()
			if opts.Color {
				switch f.Severity {
				case core.SeverityCritical:
					severity = fmt.Sprintf("\033[1;31m%s\033[0m", severity)
				case core.SeverityHigh:
					severity = fmt.Sprintf("\033[31m%s\033[0m", severity)
				case core.SeverityMedium:
					severity = fmt.Sprintf("\033[93m%s\033[0m", severity)
				case core.SeverityLow:
					severity = fmt.Sprintf("\033[33m%s\033[0m", severity)
				}
			}

			sb.WriteString(fmt.Sprintf("  [%s] %s\n", severity, f.Title))
			if f.File != "" {
				sb.WriteString(fmt.Sprintf("         File: %s:%d\n", f.File, f.Line))
			}
			sb.WriteString(fmt.Sprintf("         Rule: %s\n", f.RuleID))
			if f.Description != "" {
				sb.WriteString(fmt.Sprintf("         %s\n", f.Description))
			}
			sb.WriteString("\n")
		}
	}

	return []byte(sb.String()), nil
}
