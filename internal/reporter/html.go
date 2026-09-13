package reporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/Pixelcity-dev/Deepsec/internal/core"
)

type HTMLReporter struct{}

func (r *HTMLReporter) Name() string {
	return "html"
}

func (r *HTMLReporter) Extension() string {
	return "html"
}

func (r *HTMLReporter) Generate(results []core.ScanResult, opts ReportOptions) ([]byte, error) {
	var sb strings.Builder

	totalFindings := 0
	for _, result := range results {
		totalFindings += len(result.Findings)
	}

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>DeepSec Security Report</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #0d1117; color: #c9d1d9; padding: 20px; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: linear-gradient(135deg, #1f6feb, #388bfd); padding: 30px; border-radius: 10px; margin-bottom: 30px; }
        .header h1 { color: white; font-size: 28px; margin-bottom: 10px; }
        .header p { color: rgba(255,255,255,0.8); }
        .summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .card { background: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 20px; }
        .card h3 { color: #8b949e; font-size: 14px; margin-bottom: 10px; }
        .card .value { font-size: 32px; font-weight: bold; }
        .critical { color: #f85149; }
        .high { color: #f0883e; }
        .medium { color: #d29922; }
        .low { color: #3fb950; }
        .info { color: #58a6ff; }
        .findings { margin-top: 30px; }
        .finding { background: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 20px; margin-bottom: 15px; }
        .finding-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
        .finding-title { font-size: 18px; font-weight: 600; }
        .severity-badge { padding: 4px 12px; border-radius: 20px; font-size: 12px; font-weight: bold; text-transform: uppercase; }
        .severity-critical { background: #f85149; color: white; }
        .severity-high { background: #f0883e; color: white; }
        .severity-medium { background: #d29922; color: black; }
        .severity-low { background: #3fb950; color: black; }
        .severity-info { background: #58a6ff; color: black; }
        .finding-meta { color: #8b949e; font-size: 14px; margin-bottom: 10px; }
        .finding-desc { color: #c9d1d9; }
        .finding-code { background: #0d1117; padding: 10px; border-radius: 5px; margin-top: 10px; font-family: monospace; overflow-x: auto; }
        .footer { margin-top: 40px; text-align: center; color: #484f58; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔒 DeepSec Security Report</h1>
            <p>Generated on ` + time.Now().Format("January 2, 2006 at 15:04:05") + `</p>
        </div>
`)

	sb.WriteString(fmt.Sprintf(`
        <div class="summary">
            <div class="card">
                <h3>Total Findings</h3>
                <div class="value">%d</div>
            </div>
`, totalFindings))

	bySeverity := make(map[string]int)
	for _, result := range results {
		for _, f := range result.Findings {
			bySeverity[f.Severity.String()]++
		}
	}

	for _, sev := range []string{"CRITICAL", "HIGH", "MEDIUM", "LOW", "INFO"} {
		if count, ok := bySeverity[sev]; ok {
			class := strings.ToLower(sev)
			sb.WriteString(fmt.Sprintf(`
            <div class="card">
                <h3>%s</h3>
                <div class="value %s">%d</div>
            </div>
`, sev, class, count))
		}
	}

	sb.WriteString("        </div>\n")

	sb.WriteString(`        <div class="findings">
            <h2 style="margin-bottom: 20px;">Findings</h2>
`)

	for _, result := range results {
		for _, f := range result.Findings {
			class := strings.ToLower(f.Severity.String())
			sb.WriteString(fmt.Sprintf(`
            <div class="finding">
                <div class="finding-header">
                    <span class="finding-title">%s</span>
                    <span class="severity-badge severity-%s">%s</span>
                </div>
                <div class="finding-meta">
                    Rule: %s | File: %s:%d
                </div>
                <div class="finding-desc">%s</div>
`, f.Title, class, f.Severity.String(), f.RuleID, f.File, f.Line, f.Description))

			if f.Code != "" {
				sb.WriteString(fmt.Sprintf(`
                <div class="finding-code">%s</div>
`, f.Code))
			}

			if f.Fix != "" {
				sb.WriteString(fmt.Sprintf(`
                <div style="margin-top: 10px; padding: 10px; background: #0d4429; border-radius: 5px;">
                    <strong>Fix:</strong> %s
                </div>
`, f.Fix))
			}

			sb.WriteString("            </div>\n")
		}
	}

	sb.WriteString(`        </div>
        <div class="footer">
            <p>Generated by DeepSec v1.0.0 | https://deepsec.dev</p>
        </div>
    </div>
</body>
</html>
`)

	return []byte(sb.String()), nil
}
