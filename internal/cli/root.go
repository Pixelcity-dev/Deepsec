package cli

import (
	"fmt"
	"os"

	"github.com/Pixelcity-dev/Deepsec/internal/config"
	"github.com/Pixelcity-dev/Deepsec/internal/core"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/container"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/dast"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/format"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/iac"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/license"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/network"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/sast"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/sca"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/secrets"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/webscan"
	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	cfg      *config.Config
	registry *core.ScannerRegistry
)

var (
	version   = "1.1.0"
	buildTime = "unknown"
	commit    = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "deepsec",
	Short: "DeepSec — Cyber Security Enterprise Tool",
	Long: `DeepSec — Cyber Security Enterprise Tool

All-in-one security platform. One binary, zero dependencies, covers code + supply chain + cloud + web + formatting.

CAPABILITIES
  Code & Supply Chain  SAST (11+ langs: Go, Java, JS/TS, Python, Rust, PHP, Ruby, C/C++, C#, Kotlin, Swift)
                      SCA (11+ ecosystems: npm, pip, Maven, Go, Cargo, Composer, NuGet, etc.)
                      Secrets (200+ patterns, entropy + verified), License (SPDX/CycloneDX)
  Cloud & Containers   IaC (Terraform, CloudFormation, K8s, Dockerfile), Container (image & Dockerfile), Network
  Web & API            DAST + WebScan (OWASP Top 10, 20+ deep checks: headers, TLS, CORS, CSP, auth, etc.)
  Quality              Format (7 checks: trailing ws, EOF newline, CRLF, mixed indent, long lines, gofmt, blanks) • deepsec fmt

EXAMPLES
  deepsec scan .                                      # scan current repo
  deepsec scan . --severity high --exit-code           # gate on high/critical
  deepsec fmt . --check                               # formatting gate (CI)
  deepsec fmt . --fix                                 # auto-fix formatting
  deepsec webscan https://example.com --format sarif --output ws.sarif
  deepsec scan https://example.com --scanner webscan,dast --compliance soc2
  deepsec scan ./app --format html --output report.html

Learn more: https://pixelcity.top/deepsec  •  Docs: https://pixelcity.top/docs/deepsec
Support: service@pixelcity.dev  •  MCP: deepsec mcp start for AI agents`,
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
		initScanners()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: .deepsec.yaml, $HOME/.config/deepsec/config.yaml)")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable ANSI colors (CI)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose diagnostics")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "quiet — errors only")
	rootCmd.PersistentFlags().String("profile", "", "profile (enterprise) — overrides scanners & thresholds")
	rootCmd.PersistentFlags().String("compliance", "", "compliance mapping: soc2, iso27001, gdpr, hipaa, owasp")
	rootCmd.SetVersionTemplate(`DeepSec {{.Version}} ({{.Name}}) — Cyber Security Enterprise Tool
  commit: ` + commit + `
  built:  ` + buildTime + `
  scanners: 9  •  langs: 11+  •  https://pixelcity.top/docs/deepsec
`)
}

func initConfig() {
	if cfgFile != "" {
		var err error
		cfg, err = config.LoadConfig(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
	} else {
		cfgFile = config.FindConfigFile()
		if cfgFile != "" {
			var err error
			cfg, err = config.LoadConfig(cfgFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}
		} else {
			cfg = config.DefaultConfig()
		}
	}
}

func initScanners() {
	registry = core.NewScannerRegistry()

	if cfg.Scanners.SAST {
		registry.Register(sast.NewSASTScanner())
	}
	if cfg.Scanners.SCA {
		registry.Register(sca.NewSCAScanner())
	}
	if cfg.Scanners.Secrets {
		registry.Register(secrets.NewSecretsScanner())
	}
	if cfg.Scanners.IAC {
		registry.Register(iac.NewIACScanner())
	}
	if cfg.Scanners.Container {
		registry.Register(container.NewContainerScanner())
	}
	if cfg.Scanners.DAST {
		registry.Register(dast.NewDASTScanner())
	}
	if cfg.Scanners.Network {
		registry.Register(network.NewNetworkScanner())
	}
	if cfg.Scanners.License {
		registry.Register(license.NewLicenseScanner())
	}
	if cfg.Scanners.Format {
		registry.Register(format.NewFormatScanner())
	}
	if cfg.Scanners.WebScan {
		registry.Register(webscan.NewWebScanScanner())
	} else {
		// always register webscan even if disabled in config, so --scanner webscan works
		registry.Register(webscan.NewWebScanScanner())
	}
	// Ensure format is always available for --scanner format even if disabled
	if _, ok := registry.Get(core.ScanTypeFormat); !ok {
		registry.Register(format.NewFormatScanner())
	}
}
