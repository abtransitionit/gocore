// File: gocore/phase/phase.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com
*/
package phase

import (
	"context"
	"fmt"
	"time"

	"github.com/abtransitionit/gocore/logx"
)

// Name: Run
//
// Description:
// Run executes the function associated with the phase and returns its output and any errors.
//
// Parameters:
// - ctx: The context for the phase. This allows for cancellation and timeouts.
// - cmd: A variadic list of strings to be passed to the phase function.
//
// Returns:
// The string output and an error, if any.
func (p Phase) Run(ctx context.Context, cmd ...string) (string, error) {
	return p.fn(ctx, cmd...)
}

// Name: Run
//
// Description:
// Run executes each phase in the PhaseList sequentially.
// It logs the start and end of each phase and handles errors.
// It now also accepts a context for handling timeouts and cancellations.
//
// Parameters:
// - l: A logger instance to be used for logging progress and errors.
//
// Returns:
// An error if any phase fails or the context is canceled.
func (pl PhaseList) Run(l logx.Logger) error {
	// Create a context with a 15-minute timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel() // Always cancel the context to avoid leaks

	// Listen for OS signals to cancel the context (e.g., Ctrl+C).
	// This is important for graceful shutdown.
	go func() {
		// Example: you would listen for signals here
		// sigs := make(chan os.Signal, 1)
		// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		// <-sigs
		// l.Info("Interruption signal received, canceling workflow...")
		// cancel()
	}()

	l.Info("=== Attempting to play the sequence of phases ===")
	for _, p := range pl {
		l.Info("Starting phase '%s'...", p.Name)

		// Pass the context to the phase's Run method.
		output, err := p.Run(ctx)
		if err != nil {
			return fmt.Errorf("workflow failed on phase '%s' with error: %w", p.Name, err)
		}
		l.Info("Phase '%s' returned: %s\n\n", p.Name, output)
	}

	l.Info("=== Played the sequence of phases ===")
	return nil
}
