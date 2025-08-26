// file: gocore/cli/cli.go
package cli

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/spf13/cobra"
)

// The name of the logger to be stored in the context.
const loggerKey = "logger"

// NewCommand creates a new Cobra command with integrated gocore functionality.
// It automatically sets up logging and a consistent error handler.
func NewCommand(use, short, long string, runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		RunE:  runE,
	}
	// Add a placeholder PersistentPreRunE function to set up the context.
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Create a logger and add it to the command's context.
		logger := logx.GetLogger()
		ctx := context.WithValue(cmd.Context(), loggerKey, logger)
		cmd.SetContext(ctx)
		return nil
	}
	return cmd
}

// GetLogger retrieves the logger from the command's context.
func GetLogger(ctx context.Context) logx.Logger {
	if logger, ok := ctx.Value(loggerKey).(logx.Logger); ok {
		return logger
	}
	// Fallback to the global logger if not found.
	return logx.GetLogger()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// It is the entry point for your CLI.
func Execute(rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
