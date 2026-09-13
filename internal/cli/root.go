package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Pixelcity-dev/Deepsec/internal/config"
	"github.com/Pixelcity-dev/Deepsec/internal/core"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/container"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/dast"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/iac"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/license"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/network"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/sast"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/sca"
	"github.com/Pixelcity-dev/Deepsec/internal/scanners/secrets"
)

var (
	cfgFile  string
	cfg      *config.Config
	registry *core.ScannerRegistry
)

var rootCmd = &cobra.Command{
	Use:   "deepsec",
	Short: "DeepSec - All-in-one Cybersecurity CLI Tool",
	Long: `DeepSec is a comprehensive security scanning tool for developers and enterprises.
It provides SAST, SCA, secrets detection, IaC scanning, container security,
DAST, network scanning, and license compliance in a single binary.`,
	Version: "1.0.0",
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: .deepsec.yaml)")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "quiet mode - only output errors")
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
}
