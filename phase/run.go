// File: gocore/phase/phase.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com
*/
package phase

import (
	"context"
	"log"
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
// func (p Phase) Run(ctx context.Context, cmd ...string) (string, error) {
// 	return p.fn(ctx, cmd...)
// }

// Name: Execute
//
// Description: a placeholder method that will be implemented later.
//
// Parameters:
//
//   - ctx: The context for the workflow. This allows for cancellation and timeouts.
//
// Returns:
//
//   - An error if the workflow fails to execute.
//
// Notes:
//
//   - For now, it simply logs a message indicating that the workflow is starting.
//   - will execute each phase in the workflow sequentially.
func (w *Workflow) Execute(ctx context.Context) error {
	log.Println("Starting workflow execution...")
	// TODO: Implement topological sort and parallel execution logic here.
	log.Println("Workflow execution finished.")
	return nil
}

// Name: Run
//
// Description: executes each phase in the PhaseList sequentially.
//
// Parameters:
// - l: A logger instance to be used for logging progress and errors.
//
// Returns:
//
// An error if any phase fails or the context is canceled.
//
// Notes:
//
// - It uses a context for handling timeouts and cancellations.
// - It logs the start and end of each phase and handles errors.
// - It now also accepts a context for handling timeouts and cancellations.
// func (pl PhaseList) Run(l logx.Logger) error {
// 	// Create a context with a 15-minute timeout.
// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
// 	defer cancel() // Always cancel the context to avoid leaks

// 	// Listen for OS signals to cancel the context (e.g., Ctrl+C).
// 	// This is important for graceful shutdown.
// 	go func() {
// 		// this goroutine runs concurrently to handles the external event (Ctrl+C),
// 		// the context is the communication channel that the goroutine uses to tell the main workflow to stop.
// 		// The goroutine handles the signal, and the context propagates that signal to all the phases.
// 		// Without the context, the goroutine wouldn't have a clean, standard way to tell the rest of the application to stop its work.
// 		// Example: you would listen for signals here
// 		// sigs := make(chan os.Signal, 1)
// 		// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
// 		// <-sigs
// 		// l.Info("Interruption signal received, canceling workflow...")
// 		// cancel()
// 	}()

// 	l.Info("=== Attempting to play the sequence of phases ===")
// 	for _, p := range pl {
// 		l.Info("Starting phase '%s'...", p.Name)

// 		// Pass the context to the phase's Run method.
// 		output, err := p.Run(ctx)
// 		if err != nil {
// 			return fmt.Errorf("workflow failed on phase '%s' with error: %w", p.Name, err)
// 		}
// 		l.Info("Phase '%s' returned: %s\n\n", p.Name, output)
// 	}

// 	l.Info("=== Played the sequence of phases ===")
// 	return nil
// }
