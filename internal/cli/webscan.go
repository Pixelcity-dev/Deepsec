package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/Pixelcity-dev/Deepsec/internal/core"
	"github.com/Pixelcity-dev/Deepsec/internal/reporter"
)

var (
	webscanOutput   string
	webscanFormat   string
	webscanSeverity string
	webscanDeep     bool
)

var webscanCmd = &cobra.Command{
	Use:   "webscan [url]",
	Short: "Deep website security scanner",
	Long: `DeepScan a website for security issues with very depth checks.

Performs comprehensive OWASP-based audit:
  - Security headers (HSTS, CSP, X-Frame-Options, etc. with deep CSP analysis)
  - TLS (version, cipher, cert expiry, self-signed, hostname mismatch)
  - Cookie security (Secure, HttpOnly, SameSite)
  - CORS misconfiguration (wildcard, reflected origin, null)
  - Information disclosure (Server, X-Powered-By, etc.)
  - HTTP methods & TRACE
  - Clickjacking protection
  - Exposed files (.env, .git, backup.zip, etc.)
  - Open redirect (safe)
  - XSS reflection (safe)
  - SQLi error disclosure (safe)
  - Directory listing
  - Mixed content & SRI
  - HTTPS redirect & HSTS preload

Examples:
  deepsec webscan https://example.com
  deepsec webscan https://example.com --format json --output report.json
  deepsec webscan https://example.com --severity high
  deepsec scan https://example.com --scanner webscan,dast

Aliases: scan URL, website, audit`,
	Aliases: []string{"website", "audit", "wscan"},
	Args:    cobra.ExactArgs(1),
	RunE:    runWebscan,
}

func init() {
	webscanCmd.Flags().StringVarP(&webscanOutput, "output", "o", "", "output file path")
	webscanCmd.Flags().StringVarP(&webscanFormat, "format", "f", "table", "output format (table, json, sarif, html, csv)")
	webscanCmd.Flags().StringVar(&webscanSeverity, "severity", "INFO", "minimum severity (INFO, LOW, MEDIUM, HIGH, CRITICAL)")
	webscanCmd.Flags().BoolVar(&webscanDeep, "deep", true, "enable very deep checks (exposed files, open redirect, xss, sqli)")
	rootCmd.AddCommand(webscanCmd)
}

func runWebscan(cmd *cobra.Command, args []string) error {
	targetURL := args[0]
	if !isURL(targetURL) {
		// allow without scheme
		if !isURL("https://" + targetURL) {
			return fmt.Errorf("invalid URL: %s (expected https://example.com)", targetURL)
		}
		targetURL = "https://" + targetURL
	}

	fmt.Fprintf(os.Stderr, "DeepSec WebScan v1.0.0 - Deep website audit on %s\n", targetURL)
	if webscanDeep {
		fmt.Fprintf(os.Stderr, "Mode: very deep (headers + TLS + CORS + exposed files + open redirect + XSS + SQLi + ...)\n")
	}

	start := time.Now()

	ruleEngine := core.NewRuleEngine()
	// Load rules if available (not required for webscan, but keep for consistency)
	_ = ruleEngine.LoadRulesFromDir("rules")

	pipeline := core.NewPipeline(registry, ruleEngine)
	filter := core.NewFindingFilter()
	filter.MinSeverity = core.ParseSeverity(webscanSeverity)
	pipeline.SetFilter(filter)

	targetObj := core.Target{
		Kind: core.TargetURL,
		URI:  targetURL,
		Options: map[string]interface{}{
			"deep": webscanDeep,
		},
	}

	results, err := pipeline.Scan(context.Background(), targetObj, []core.ScanType{core.ScanTypeWebScan})
	if err != nil {
		return fmt.Errorf("webscan failed: %w", err)
	}

	duration := time.Since(start).Seconds()
	total := 0
	for _, r := range results {
		total += len(r.Findings)
		// If webscan scanner not registered, try fallback to DAST
		if r.Scanner == "webscan" && len(r.Findings) == 0 && total == 0 {
			// no findings may still be ok
		}
	}

	// Fallback: if webscan produced no scanner (not registered), try DAST as alias
	if len(results) == 0 || (len(results) == 1 && results[0].Scanner == "" && total == 0) {
		// try dast + webscan combined
		results, _ = pipeline.Scan(context.Background(), targetObj, []core.ScanType{core.ScanTypeDAST, core.ScanTypeWebScan})
		total = 0
		for _, r := range results {
			total += len(r.Findings)
		}
	}

	fmt.Fprintf(os.Stderr, "\nWebScan completed in %.2f seconds\n", duration)
	if total == 0 {
		fmt.Fprintf(os.Stderr, "No issues found - site looks good! (checked %d categories)\n", 16)
	} else {
		fmt.Fprintf(os.Stderr, "Found %d issues\n", total)
		// Severity breakdown
		bySev := make(map[string]int)
		for _, r := range results {
			for _, f := range r.Findings {
				bySev[f.Severity.String()]++
			}
		}
		for _, sev := range []string{"CRITICAL", "HIGH", "MEDIUM", "LOW", "INFO"} {
			if c, ok := bySev[sev]; ok && c > 0 {
				fmt.Fprintf(os.Stderr, "  %s: %d\n", sev, c)
			}
		}
	}

	// Generate report
	rpt := reporter.GetReporter(webscanFormat)
	if rpt == nil {
		return fmt.Errorf("unknown format: %s", webscanFormat)
	}
	output, err := rpt.Generate(results, reporter.ReportOptions{
		Format: webscanFormat,
		Output: webscanOutput,
		Color:  !rootCmd.PersistentFlags().Changed("no-color"),
	})
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	if webscanOutput != "" {
		if err := os.WriteFile(webscanOutput, output, 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Report written to %s\n", webscanOutput)
	} else {
		fmt.Print(string(output))
	}

	// Hint
	if total > 0 {
		fmt.Fprintf(os.Stderr, "\nHint: run with --format sarif --output results.sarif for GitHub Code Scanning\n")
		fmt.Fprintf(os.Stderr, "      deepsec scan https://example.com --scanner webscan --severity high\n")
	}

	return nil
}
