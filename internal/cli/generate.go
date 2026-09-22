package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate configuration files for CI/CD",
	Long:  `Generate configuration files for various CI/CD platforms and tools.`,
}

var generatePreCommitCmd = &cobra.Command{
	Use:   "pre-commit",
	Short: "Generate .pre-commit-config.yaml",
	RunE:  runGeneratePreCommit,
}

var generateGitHubCmd = &cobra.Command{
	Use:   "github-action",
	Short: "Generate GitHub Actions workflow",
	RunE:  runGenerateGitHub,
}

var generateGitLabCmd = &cobra.Command{
	Use:   "gitlab-ci",
	Short: "Generate GitLab CI configuration",
	RunE:  runGenerateGitLabCI,
}

func init() {
	generateCmd.AddCommand(generatePreCommitCmd)
	generateCmd.AddCommand(generateGitHubCmd)
	generateCmd.AddCommand(generateGitLabCmd)
	rootCmd.AddCommand(generateCmd)
}

func runGeneratePreCommit(cmd *cobra.Command, args []string) error {
	content := `repos:
  - repo: https://github.com/deepsec/deepsec
    rev: v1.1.0
    hooks:
      - id: deepsec-secrets
        name: DeepSec Secrets Scan
      - id: deepsec-sast
        name: DeepSec SAST Scan
      - id: deepsec-iac
        name: DeepSec IaC Scan
`

	if err := os.WriteFile(".pre-commit-config.yaml", []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write pre-commit config: %w", err)
	}

	fmt.Println("Created .pre-commit-config.yaml")
	return nil
}

func runGenerateGitHub(cmd *cobra.Command, args []string) error {
	dir := filepath.Join(".github", "workflows")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	content := `name: DeepSec Security Scan

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]
  schedule:
    - cron: '0 0 * * 0'

permissions:
  contents: read
  security-events: write

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Run DeepSec SAST
        uses: deepsec/deepsec-action@v1
        with:
          scan-type: 'sast'
          format: 'sarif'
          output: 'sast-results.sarif'

      - name: Upload SAST results to GitHub Code Scanning
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'sast-results.sarif'
        if: always()

      - name: Run DeepSec SCA
        uses: deepsec/deepsec-action@v1
        with:
          scan-type: 'sca'
          format: 'sarif'
          output: 'sca-results.sarif'

      - name: Upload SCA results to GitHub Code Scanning
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'sca-results.sarif'
        if: always()

      - name: Run DeepSec Secrets
        uses: deepsec/deepsec-action@v1
        with:
          scan-type: 'secrets'
          format: 'sarif'
          output: 'secrets-results.sarif'

      - name: Upload Secrets results to GitHub Code Scanning
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'secrets-results.sarif'
        if: always()
`

	path := filepath.Join(dir, "deepsec.yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write GitHub Actions config: %w", err)
	}

	fmt.Printf("Created %s\n", path)
	return nil
}

func runGenerateGitLabCI(cmd *cobra.Command, args []string) error {
	content := `stages:
  - security

deepsec-sast:
  stage: security
  image: deepsec/deepsec:latest
  script:
    - deepsec scan fs --format sarif --output gl-sast-report.json --scanner sast
  artifacts:
    reports:
      sast: gl-sast-report.json
  only:
    - branches

deepsec-sca:
  stage: security
  image: deepsec/deepsec:latest
  script:
    - deepsec scan fs --format cyclonedx --output gl-dependency-report.json --scanner sca
  artifacts:
    reports:
      dependency_scanning: gl-dependency-report.json
  only:
    - branches

deepsec-secrets:
  stage: security
  image: deepsec/deepsec:latest
  script:
    - deepsec scan fs --format json --output gl-secret-detection-report.json --scanner secrets
  artifacts:
    reports:
      secret_detection: gl-secret-detection-report.json
  only:
    - branches
`

	if err := os.WriteFile(".gitlab-ci.yml", []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write GitLab CI config: %w", err)
	}

	fmt.Println("Created .gitlab-ci.yml")
	return nil
}
