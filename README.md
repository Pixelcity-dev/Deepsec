# DeepSec

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Pixelcity-dev/Deepsec)](https://goreportcard.com/report/github.com/Pixelcity-dev/Deepsec)
[![Release](https://img.shields.io/github/v/release/Pixelcity-dev/Deepsec)](https://github.com/Pixelcity-dev/Deepsec/releases)

**All-in-one Cybersecurity CLI Tool for Developers and Enterprise**

DeepSec is a comprehensive security scanning tool that provides SAST, SCA, secrets detection, IaC scanning, container security, DAST, network scanning, and license compliance in a single, fast, zero-dependency binary.

## Features

- **SAST** - Static Application Security Testing for 11+ languages
- **SCA** - Software Composition Analysis for 11+ ecosystems
- **Secrets** - Detect leaked credentials and API keys
- **IaC** - Infrastructure as Code security scanning
- **Container** - Container image and Dockerfile analysis
- **DAST** - Dynamic Application Security Testing
- **Network** - Network vulnerability scanning
- **License** - License compliance checking

## Installation

### Binary Download

```bash
curl -sSL https://pixelcity.dev/deepsec/install.sh | sh
```

### Homebrew

```bash
brew install deepsec
```

### Docker

```bash
docker pull deepsec/deepsec:latest
```

### Build from Source

```bash
git clone https://github.com/Pixelcity-dev/Deepsec.git
cd Deepsec
make build
```

## Quick Start

```bash
# Scan current directory
deepsec scan .

# Scan with specific scanner
deepsec scan . --scanner sast,sca,secrets

# Scan with severity filter
deepsec scan . --severity high,critical

# Output to JSON
deepsec scan . --format json --output results.json

# Output to SARIF (for GitHub Code Scanning)
deepsec scan . --format sarif --output results.sarif

# Initialize configuration
deepsec init

# Update vulnerability databases
deepsec db update
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `deepsec scan [target]` | Run security scan |
| `deepsec init` | Initialize configuration |
| `deepsec db update` | Update vulnerability databases |
| `deepsec rule list` | List security rules |
| `deepsec rule search [query]` | Search rules |
| `deepsec plugin list` | List installed plugins |
| `deepsec plugin install [name]` | Install a plugin |
| `deepsec server start` | Start scanning server |
| `deepsec mcp start` | Start MCP server for AI agents |
| `deepsec convert` | Convert between formats |
| `deepsec generate` | Generate CI/CD configs |

## Configuration

Create a `.deepsec.yaml` file in your project root:

```yaml
version: "1.0.0"
scanners:
  sast: true
  sca: true
  secrets: true
  iac: true
  container: true
  dast: false
  network: false
  license: true
report:
  format: table
  color: true
filter:
  min_severity: low
```

## CI/CD Integration

### GitHub Actions

```yaml
- uses: deepsec/deepsec-action@v1
  with:
    scan-type: 'sast,sca,secrets'
    severity: 'high,critical'
    format: 'sarif'
    upload-code-scanning: true
```

### GitLab CI

```yaml
deepsec-sast:
  image: deepsec/deepsec:latest
  script:
    - deepsec scan fs --format sarif --output results.sarif
  artifacts:
    reports:
      sast: results.sarif
```

### Pre-commit Hook

```yaml
repos:
  - repo: https://github.com/deepsec/deepsec
    rev: v1.0.0
    hooks:
      - id: deepsec-secrets
      - id: deepsec-sast
```

## Output Formats

| Format | Extension | Use Case |
|--------|-----------|----------|
| Table | (terminal) | Human-readable console output |
| JSON | `.json` | Machine-readable, pipeline integration |
| SARIF 2.1.0 | `.sarif` | GitHub Code Scanning, SonarQube |
| CycloneDX | `.json` | SBOM standard, dependency tracking |
| SPDX | `.json` | SBOM alternative, license compliance |
| HTML | `.html` | Visual reports, stakeholder sharing |
| JUnit XML | `.xml` | CI/CD test result integration |
| CSV | `.csv` | Spreadsheet analysis |

## Plugin System

```bash
# Search for plugins
deepsec plugin search sast

# Install a plugin
deepsec plugin install custom-scanner

# List installed plugins
deepsec plugin list
```

## MCP Server (AI Agents)

```bash
# Start MCP server for AI coding assistants
deepsec mcp start
```

Exposed tools:
- `deepsec.scan` - Trigger a security scan
- `deepsec.findings` - Read scan findings
- `deepsec.explain` - Get explanation of a finding
- `deepsec.suggest-fix` - Get remediation suggestions

## Documentation

- [User Guide](docs/user-guide.md)
- [Rule Authoring](docs/rules.md)
- [Plugin Development](docs/plugins.md)
- [API Reference](docs/api.md)

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) first.

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
