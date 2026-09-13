package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Pixelcity-dev/Deepsec/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize DeepSec configuration",
	Long:  `Create a default .deepsec.yaml configuration file in the current directory.`,
	RunE:  runInit,
}

func init() {
	initCmd.Flags().Bool("force", false, "overwrite existing config file")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	cfg := config.DefaultConfig()

	force, _ := cmd.Flags().GetBool("force")
	if !force {
		if _, err := os.Stat(".deepsec.yaml"); err == nil {
			return fmt.Errorf("config file already exists. Use --force to overwrite")
		}
	}

	if err := cfg.Save(".deepsec.yaml"); err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	fmt.Println("Created .deepsec.yaml")
	return nil
}
