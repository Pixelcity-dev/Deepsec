package dast

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Pixelcity-dev/Deepsec/internal/core"
)

type DASTScanner struct {
	core.BaseScanner
	name string
}

func NewDASTScanner() *DASTScanner {
	return &DASTScanner{
		BaseScanner: core.BaseScanner{Enabled: true},
		name:        "dast",
	}
}

func (s *DASTScanner) Name() string {
	return s.name
}

func (s *DASTScanner) Type() core.ScanType {
	return core.ScanTypeDAST
}

func (s *DASTScanner) SupportedTargets() []core.TargetKind {
	return []core.TargetKind{core.TargetURL}
}

func (s *DASTScanner) Scan(ctx context.Context, target core.Target, rules []core.Rule) ([]core.Finding, error) {
	var findings []core.Finding

	targetURL, err := url.Parse(target.URI)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %w", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	findings = append(findings, checkSecurityHeaders(targetURL, client)...)
	findings = append(findings, checkTLSConfiguration(targetURL)...)
	findings = append(findings, checkCookieSecurity(targetURL, client)...)

	return findings, nil
}

func checkSecurityHeaders(targetURL *url.URL, client *http.Client) []core.Finding {
	var findings []core.Finding

	resp, err := client.Get(targetURL.String())
	if err != nil {
		return findings
	}
	defer resp.Body.Close()

	requiredHeaders := map[string]struct {
		Name     string
		Severity core.Severity
		Fix      string
	}{
		"X-Content-Type-Options": {
			Name:     "X-Content-Type-Options",
			Severity: core.SeverityMedium,
			Fix:      "Add header: X-Content-Type-Options: nosniff",
		},
		"X-Frame-Options": {
			Name:     "X-Frame-Options",
			Severity: core.SeverityMedium,
			Fix:      "Add header: X-Frame-Options: DENY or SAMEORIGIN",
		},
		"X-XSS-Protection": {
			Name:     "X-XSS-Protection",
			Severity: core.SeverityLow,
			Fix:      "Add header: X-XSS-Protection: 1; mode=block",
		},
		"Strict-Transport-Security": {
			Name:     "Strict-Transport-Security",
			Severity: core.SeverityHigh,
			Fix:      "Add header: Strict-Transport-Security: max-age=31536000; includeSubDomains",
		},
		"Content-Security-Policy": {
			Name:     "Content-Security-Policy",
			Severity: core.SeverityHigh,
			Fix:      "Add Content-Security-Policy header with appropriate directives",
		},
		"Referrer-Policy": {
			Name:     "Referrer-Policy",
			Severity: core.SeverityMedium,
			Fix:      "Add header: Referrer-Policy: strict-origin-when-cross-origin",
		},
		"Permissions-Policy": {
			Name:     "Permissions-Policy",
			Severity: core.SeverityMedium,
			Fix:      "Add Permissions-Policy header to restrict browser features",
		},
	}

	for headerName, info := range requiredHeaders {
		if resp.Header.Get(headerName) == "" {
			findings = append(findings, core.Finding{
				RuleID:      "dast-missing-header-" + strings.ToLower(headerName),
				Severity:    info.Severity,
				Category:    "security-headers",
				Title:       "Missing " + info.Name + " header",
				Description: "Response does not include " + info.Name + " header",
				Fix:         info.Fix,
				References: []string{
					"https://securityheaders.com/",
					"https://owasp.org/www-project-secure-headers/",
				},
				Confidence: 1.0,
			})
		}
	}

	serverHeader := resp.Header.Get("Server")
	if serverHeader != "" {
		findings = append(findings, core.Finding{
			RuleID:      "dast-server-header-disclosure",
			Severity:    core.SeverityLow,
			Category:    "information-disclosure",
			Title:       "Server header disclosure",
			Description: "Server header reveals: " + serverHeader,
			Fix:         "Remove or obfuscate Server header",
			Confidence:  0.9,
		})
	}

	return findings
}

func checkTLSConfiguration(targetURL *url.URL) []core.Finding {
	var findings []core.Finding

	if targetURL.Scheme == "http" {
		findings = append(findings, core.Finding{
			RuleID:      "dast-no-tls",
			Severity:    core.SeverityHigh,
			Category:    "transport-security",
			Title:       "HTTP instead of HTTPS",
			Description: "Target uses HTTP instead of HTTPS",
			Fix:         "Use HTTPS and implement HTTP to HTTPS redirect",
			References: []string{"https://owasp.org/www-project-transport-layer-security-cheat-sheet/"},
			Confidence:  1.0,
		})
	}

	return findings
}

func checkCookieSecurity(targetURL *url.URL, client *http.Client) []core.Finding {
	var findings []core.Finding

	resp, err := client.Get(targetURL.String())
	if err != nil {
		return findings
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if !cookie.Secure {
			findings = append(findings, core.Finding{
				RuleID:      "dast-cookie-no-secure",
				Severity:    core.SeverityMedium,
				Category:    "cookie-security",
				Title:       "Cookie without Secure flag: " + cookie.Name,
				Description: "Cookie '" + cookie.Name + "' does not have Secure flag",
				Fix:         "Set Secure flag on cookie: " + cookie.Name,
				Confidence:  1.0,
			})
		}

		if !cookie.HttpOnly {
			findings = append(findings, core.Finding{
				RuleID:      "dast-cookie-no-httponly",
				Severity:    core.SeverityMedium,
				Category:    "cookie-security",
				Title:       "Cookie without HttpOnly flag: " + cookie.Name,
				Description: "Cookie '" + cookie.Name + "' does not have HttpOnly flag",
				Fix:         "Set HttpOnly flag on cookie: " + cookie.Name,
				Confidence:  1.0,
			})
		}
	}

	return findings
}

func init() {
	_ = fmt.Sprintf
}
