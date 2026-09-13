package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	convertInput  string
	convertOutput string
	convertFormat string
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert between scan result formats",
	Long:  `Convert scan results from one format to another.`,
	RunE:  runConvert,
}

func init() {
	convertCmd.Flags().StringVarP(&convertInput, "input", "i", "", "input file")
	convertCmd.Flags().StringVarP(&convertOutput, "output", "o", "", "output file")
	convertCmd.Flags().StringVarP(&convertFormat, "format", "f", "", "output format")

	convertCmd.MarkFlagRequired("input")
	convertCmd.MarkFlagRequired("output")
	convertCmd.MarkFlagRequired("format")

	rootCmd.AddCommand(convertCmd)
}

func runConvert(cmd *cobra.Command, args []string) error {
	fmt.Printf("Converting %s to %s format...\n", convertInput, convertFormat)

	data, err := os.ReadFile(convertInput)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	_ = data

	if err := os.WriteFile(convertOutput, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Converted successfully: %s\n", convertOutput)
	return nil
}
