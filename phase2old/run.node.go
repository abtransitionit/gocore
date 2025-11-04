package phase2

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/run"
)

// CustomErrorHandler is your function type to handle non-critical errors
type CustomErrorHandler func(err error, logger logx.Logger) bool

type Node struct {
	Name       string
	IsLocal    bool
	RemoteHost string
	Cli        string
}

// Execute runs the CLI command on this node (local or remote) and handles errors
// func (node *Node) Execute(ctx context.Context, logger logx.Logger, errorHandler CustomErrorHandler) (string, error) {
func (node *Node) Execute(ctx context.Context, logger logx.Logger) (string, error) {
	logger.Infof("🅝 Node %s > Starting CLI execution", node.Name)

	output, err := run.ExecuteCliQuery(node.Cli, logger, node.IsLocal, node.RemoteHost, run.NoOpErrorHandler)
	if err != nil {
		logger.Errorf("🅝 Node %s > CLI execution failed: %v", node.Name, err)
		return output, fmt.Errorf("node %s: %w", node.Name, err)
	}

	logger.Infof("🅝 Node %s > CLI execution completed successfully", node.Name)
	return output, nil
}

// // Optional: Execute multiple nodes concurrently (like GoFunc.Execute)
// func ExecuteNodesConcurrently(ctx context.Context, nodes []*Node, logger logx.Logger, errorHandler CustomErrorHandler) error {
// 	var wg sync.WaitGroup
// 	errCh := make(chan error, len(nodes))

// 	for _, n := range nodes {
// 		wg.Add(1)
// 		node := n // capture variable
// 		go func() {
// 			defer wg.Done()
// 			if _, err := node.Execute(ctx, logger, errorHandler); err != nil {
// 				errCh <- err
// 			}
// 		}()
// 	}

// 	wg.Wait()
// 	close(errCh)

// 	var nodeErrs []error
// 	for e := range errCh {
// 		nodeErrs = append(nodeErrs, e)
// 	}

// 	if len(nodeErrs) > 0 {
// 		return fmt.Errorf("%d node(s) failed", len(nodeErrs))
// 	}

// 	return nil
// }
