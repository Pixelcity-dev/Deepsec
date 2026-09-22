# DeepSec — Cyber Security Enterprise Tool

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Pixelcity-dev/Deepsec)](https://goreportcard.com/report/github.com/Pixelcity-dev/Deepsec)
[![Release](https://img.shields.io/github/v/release/Pixelcity-dev/Deepsec)](https://github.com/Pixelcity-dev/Deepsec/releases)

**Cyber Security Enterprise Tool — code, supply chain, cloud and web in one binary.**

DeepSec unifies **SAST, SCA, Secrets, IaC, Containers, DAST, WebScan, Network, License & SBOM** for **any language** and **any webpage**. Zero dependencies, <50ms startup, compliance-ready (SOC 2, ISO 27001, GDPR, OWASP, CWE).

## Features — Any Language. Any Webpage.

- **SAST** — 11+ langs: Go, Java, JS/TS, Python, Rust, PHP, Ruby, C/C++, C#, Kotlin, Swift • OWASP A03
- **SCA** — 11+ ecosystems: npm, pip, Maven, Go, Cargo, Composer, NuGet, RubyGems • OWASP A06
- **Secrets** — 200+ patterns (AWS, GCP, Stripe…), entropy + verified • OWASP A07
- **IaC** — Terraform, CloudFormation, Kubernetes, Dockerfile (CIS) • OWASP A05
- **Container** — Image layers, Dockerfile, distroless checks
- **DAST** — Headers, TLS, CORS, cookies, clickjacking • OWASP A01/A05
- **WebScan** — 20+ deep checks (HSTS/CSP deep, TLS cert/cipher/version, CORS wildcard, exposed `/.env/.git`, open redirect, XSS reflection, SQLi error, directory listing, `security.txt`, mixed content, SRI, HTTPS redirect) • OWASP Top 10
- **Network** — Port/service, TLS, header audit
- **License & SBOM** — SPDX/CycloneDX, license compliance
- **Any Webpage** — `deepsec webscan https://example.com` audits marketing sites, SaaS apps, APIs, SPAs

## Installation

### One-Liner — CDN via Caddy

```bash
curl -fsSL https://cdn.pixelcity.dev/deepsec/install.sh | sh
# fallbacks
curl -fsSL https://cdn.pixelcity.top/deepsec/install.sh | sh
curl -fsSL https://pixelcity.top/deepsec/install.sh | sh
wget -qO- https://cdn.pixelcity.dev/deepsec/install.sh | sh
```

```bash
DEEPSEC_VERSION=latest curl -fsSL https://cdn.pixelcity.dev/deepsec/install.sh | sh
INSTALL_DIR=/usr/local/bin curl -fsSL https://cdn.pixelcity.dev/deepsec/install.sh | sh
curl -fsSL https://cdn.pixelcity.dev/deepsec/install.sh | sh -s -- --help
```

### Go Install

```bash
go install github.com/Pixelcity-dev/Deepsec/cmd/deepsec@latest
```

### Build from Source

```bash
git clone https://github.com/Pixelcity-dev/Deepsec.git
# private SSH
git clone git@github.com:Pixelcity-dev/Deepsec.git
cd Deepsec
make build
sudo make install
```

### Manual Download

```bash
# https://cdn.pixelcity.dev/deepsec/releases/v1.0.0/
curl -fsSL https://cdn.pixelcity.dev/deepsec/releases/v1.0.0/deepsec-linux-amd64 -o deepsec
chmod +x deepsec && sudo mv deepsec /usr/local/bin/
```

Releases served via `Caddy (443) → Nginx CDN (pixelcity-cdn)` from `/opt/pixelcity/cdn/assets/deepsec/` with `Cache-Control` & CORS.

## Quick Start

```bash
# Scan current directory
deepsec scan .

# Specific scanners
deepsec scan . --scanner sast,sca,secrets

# Deep Website Scan — 20+ checks
deepsec webscan https://example.com
deepsec webscan https://example.com --format json --output report.json
deepsec webscan https://example.com --severity high --format sarif --output webscan.sarif
deepsec scan https://example.com --scanner webscan,dast

# Severity filter
deepsec scan . --severity high,critical

# JSON / SARIF
deepsec scan . --format json --output results.json
deepsec scan . --format sarif --output results.sarif

# Config & DB
deepsec init
deepsec db update
```

## CLI

| Command | Description |
|---------|-------------|
| `deepsec scan [target]` | Scan `fs`/`url`/`image`/`repo` — auto-detects target |
| `deepsec webscan [url]` | Deep website audit — 20+ OWASP checks (`website`/`audit`/`wscan` aliases) |
| `deepsec init` | Scaffold `.deepsec.yaml` |
| `deepsec db update` | Update NVD/OSV DB (air-gapped cache) |
| `deepsec rule list/search` | 1000+ rules, filter by `language/category/severity` |
| `deepsec plugin list/install` | Custom scanners, private registry |
| `deepsec server start` | HTTP service (`:8443`, TLS, audit log) |
| `deepsec mcp start` | MCP for AI agents (`deepsec.scan`, `deepsec.explain`) |
| `deepsec convert/generate` | Format convert, GitHub/GitLab/pre-commit generators |

Exit codes: `0` pass, `1` gate failed, `2` error. Flags: `--severity`, `--format`, `--output`, `--profile`, `--compliance`, `--fail-on`, `--no-color`.

## Configuration

```yaml
version: "1.0.0"
scanners: { sast: true, sca: true, secrets: true, iac: true, container: true, dast: true, webscan: true, network: true, license: true }
report: { format: html, color: true }   # table|json|sarif|cyclonedx|spdx|html|junit|csv
filter:
  min_severity: low          # INFO|LOW|MEDIUM|HIGH|CRITICAL
  fail_on: HIGH              # gate for CI
  compliance: [soc2, iso27001, gdpr, owasp]
  exclude_rules: [webscan-missing-security-txt]
```

CLI overrides: `--profile enterprise --compliance soc2 --fail-on high --severity medium`

Profiles: `enterprise` preset via `--profile` or `filter.profile`.

## Website Scanner — Deep Security Check

```bash
deepsec webscan https://example.com
deepsec webscan https://example.com --format json --output webscan.json
deepsec webscan https://pixelcity.top --severity medium
deepsec scan https://example.com --scanner webscan --format sarif
```

**20+ Checks:**
- **Security Headers Deep** - HSTS (max-age, includeSubDomains, preload), CSP (unsafe-inline, wildcard), X-Frame-Options, COOP/COEP/CORP
- **TLS Deep** - TLS version, weak ciphers, cert expiry, self-signed, hostname mismatch
- **Cookie Security** - Secure, HttpOnly, SameSite
- **CORS** - Wildcard *, reflected Origin, null
- **Information Disclosure** - Server, X-Powered-By, Via
- **HTTP Methods & TRACE** - OPTIONS, TRACE/XST
- **Clickjacking** - X-Frame-Options / CSP frame-ancestors
- **Exposed Files** - `/.env`, `/.git/HEAD`, backup.zip, config.json
- **Open Redirect** - `?url=//evil.com` (safe)
- **XSS Reflection** - safe marker payload
- **SQLi Error Disclosure** - `'` injection (safe)
- **Directory Listing** - Index of/
- **Security.txt & robots.txt** - RFC 9116
- **Mixed Content & SRI** - http on https, missing integrity
- **HTTPS Redirect** - http → https

Example output (pixelcity.top):
```
DeepSec WebScan v1.0.0 - Deep website audit on https://pixelcity.top
WebScan completed in 1.36 seconds
Found 7 issues  MEDIUM:1 LOW:5 INFO:1
```

## Reporting & Compliance

Every finding mapped to **OWASP Top 10, CWE, SOC 2, ISO 27001, GDPR** (`internal/config/config.go:14`). Table shows `Risk Score` (Critical 40, High 10, Medium 3, Low 1) + `Risk Level` (Excellent/Low/Medium/High/Critical). HTML is executive report: risk meter, exposure bar, compliance grid, prioritized fix plan.

```bash
deepsec scan . --format html --output report.html
deepsec scan . --format sarif --output results.sarif
deepsec scan . --format cyclonedx --output sbom.json
deepsec webscan https://example.com --compliance soc2 --format html --output ws.html
```

## CI/CD

### GitHub Actions

```yaml
- uses: deepsec/deepsec-action@v1
  with:
    scan-type: 'sast,sca,secrets,webscan'
    severity: 'high,critical'
    compliance: 'soc2,owasp'
    format: 'sarif'
- uses: github/codeql-action/upload-sarif@v3
  with: { sarif_file: results.sarif }
```

### GitLab CI

```yaml
deepsec:
  image: deepsec/deepsec:latest
  script:
    - deepsec scan . --profile enterprise --format sarif --output gl-sast.json
    - deepsec scan . --format cyclonedx --output sbom.json
  artifacts: { reports: { sast: gl-sast.json }, paths: [sbom.json] }
```

### Pre-commit

```yaml
repos:
  - repo: https://github.com/Pixelcity-dev/Deepsec
    rev: v1.0.0
    hooks: [{id: deepsec-secrets}, {id: deepsec-sast}]
```

### Policy Gate

```bash
deepsec scan . --fail-on HIGH
deepsec scan . --fail-on MEDIUM
```

## Output Formats

| Format | Use |
|--------|-----|
| **Table** | Terminal, CI logs — risk score + compliance, grouped by severity |
| **HTML** | Executive report — risk meter, compliance grid, fix plan |
| **JSON** | Automation — `jq` friendly |
| **SARIF 2.1.0** | GitHub Code Scanning, SonarQube, SIEM |
| **CycloneDX/SPDX** | SBOM, supply-chain |
| **JUnit/CSV** | Test gates, spreadsheets |

`--format html|table|json|sarif|cyclonedx|spdx|junit|csv` + `--output`

## Plugin System

```bash
deepsec plugin search sast
deepsec plugin install custom-scanner --registry https://plugins.yourco.dev
deepsec plugin list
```

## MCP Server — AI-Native

```bash
deepsec mcp start   # stdio/SSE
```

Tools: `deepsec.scan`, `deepsec.findings`, `deepsec.explain`, `deepsec.suggest-fix`

## Deployment

```bash
# Air-gapped
deepsec db update --cache-dir /mnt/cache && tar czf deepsec-db.tar.gz ~/.deepsec/cache
# Server
deepsec server start --host 0.0.0.0 --port 8443 --tls
```

## Documentation

- [User Guide](docs/user-guide.md) • [Rule Authoring](docs/rules.md) • [Plugin Dev](docs/plugins.md) • [API](docs/api.md)
- Full Docs: **https://pixelcity.top/docs/deepsec** • **https://pixelcity.dev/docs/deepsec** • **https://cdn.pixelcity.dev/docs/deepsec**

## Contributing & License

Contributions via PR — see `CONTRIBUTING.md`. Security reports and help: `service@pixelcity.dev`.
**Apache 2.0** — `LICENSE` — `https://pixelcity.top`
