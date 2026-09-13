package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  string            `yaml:"version" json:"version"`
	Scanners ScannerConfig     `yaml:"scanners" json:"scanners"`
	Report   ReportConfig      `yaml:"report" json:"report"`
	Rules    RulesConfig       `yaml:"rules" json:"rules"`
	DB       DBConfig          `yaml:"db" json:"db"`
	Plugin   PluginConfig      `yaml:"plugin" json:"plugin"`
	Server   ServerConfig      `yaml:"server" json:"server"`
	Filter   FilterConfig      `yaml:"filter" json:"filter"`
}

type ScannerConfig struct {
	SAST      bool `yaml:"sast" json:"sast"`
	SCA       bool `yaml:"sca" json:"sca"`
	Secrets   bool `yaml:"secrets" json:"secrets"`
	IAC       bool `yaml:"iac" json:"iac"`
	Container bool `yaml:"container" json:"container"`
	DAST      bool `yaml:"dast" json:"dast"`
	Network   bool `yaml:"network" json:"network"`
	License   bool `yaml:"license" json:"license"`
}

type ReportConfig struct {
	Format    string `yaml:"format" json:"format"`
	Output    string `yaml:"output" json:"output"`
	Template  string `yaml:"template,omitempty" json:"template,omitempty"`
	Color     bool   `yaml:"color" json:"color"`
	Verbose   bool   `yaml:"verbose" json:"verbose"`
}

type RulesConfig struct {
 Paths    []string `yaml:"paths" json:"paths"`
 BuiltIn  bool     `yaml:"builtin" json:"builtin"`
 Custom   []string `yaml:"custom,omitempty" json:"custom,omitempty"`
}

type DBConfig struct {
	CacheDir  string `yaml:"cache_dir" json:"cache_dir"`
	UpdateURL string `yaml:"update_url" json:"update_url"`
	AutoUpdate bool `yaml:"auto_update" json:"auto_update"`
}

type PluginConfig struct {
	Dir      string            `yaml:"dir" json:"dir"`
	Registry string            `yaml:"registry" json:"registry"`
	Enabled  map[string]bool   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
}

type ServerConfig struct {
	Host string `yaml:"host" json:"host"`
	Port int    `yaml:"port" json:"port"`
	TLS  bool   `yaml:"tls" json:"tls"`
}

type FilterConfig struct {
	MinSeverity     string   `yaml:"min_severity" json:"min_severity"`
	ExcludeRules    []string `yaml:"exclude_rules,omitempty" json:"exclude_rules,omitempty"`
	IncludeRules    []string `yaml:"include_rules,omitempty" json:"include_rules,omitempty"`
	ExcludeFiles    []string `yaml:"exclude_files,omitempty" json:"exclude_files,omitempty"`
	ExcludeCategories []string `yaml:"exclude_categories,omitempty" json:"exclude_categories,omitempty"`
}

func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".deepsec", "cache")
	pluginDir := filepath.Join(homeDir, ".deepsec", "plugins")

	return &Config{
		Version: "1.0.0",
		Scanners: ScannerConfig{
			SAST:      true,
			SCA:       true,
			Secrets:   true,
			IAC:       true,
			Container: true,
			DAST:      true,
			Network:   true,
			License:   true,
		},
		Report: ReportConfig{
			Format: "table",
			Color:  true,
		},
		Rules: RulesConfig{
			Paths:   []string{},
			BuiltIn: true,
		},
		DB: DBConfig{
			CacheDir:   cacheDir,
			AutoUpdate: true,
		},
		Plugin: PluginConfig{
			Dir:      pluginDir,
			Registry: "https://plugins.deepsec.dev",
			Enabled:  make(map[string]bool),
		},
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8443,
		},
		Filter: FilterConfig{
			MinSeverity: "INFO",
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func FindConfigFile() string {
	candidates := []string{
		".deepsec.yaml",
		"deepsec.yaml",
		".deepsec.yml",
		"deepsec.yml",
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}

	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		homeConfig := filepath.Join(homeDir, ".config", "deepsec", "config.yaml")
		if _, err := os.Stat(homeConfig); err == nil {
			return homeConfig
		}
	}

	return ""
}

func (c *Config) MergeWith(other *Config) {
	if other.Scanners.SAST {
		c.Scanners.SAST = true
	}
	if other.Scanners.SCA {
		c.Scanners.SCA = true
	}
	if other.Scanners.Secrets {
		c.Scanners.Secrets = true
	}
	if other.Scanners.IAC {
		c.Scanners.IAC = true
	}
	if other.Scanners.Container {
		c.Scanners.Container = true
	}
	if other.Scanners.DAST {
		c.Scanners.DAST = true
	}
	if other.Scanners.Network {
		c.Scanners.Network = true
	}
	if other.Scanners.License {
		c.Scanners.License = true
	}
}
