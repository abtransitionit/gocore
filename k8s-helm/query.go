package helm

import (
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/run"
)

func QueryHelm(helmHost, helmQuery string, logger logx.Logger) (string, error) {
	var output string
	var err error

	if helmHost == "" {
		output, err = run.RunOnLocal(helmQuery)
	} else {
		output, err = run.RunCliSsh(helmHost, helmQuery)
	}

	// Handle "helm errors" that are not true errors
	if HandleHelmError(err, logger) {
		return "", nil // handled gracefully, nothing to print
	}

	// success
	return output, nil
}

func HandleHelmError(err error, logger logx.Logger) bool {
	if err == nil {
		return false
	}

	msgLower := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msgLower, "no repositories to show"):
		logger.Warnf("No repositories configured yet")
		return true

	case strings.Contains(msgLower, "release: not found"):
		logger.Warnf("No release found for this query")
		return true

	case strings.Contains(msgLower, "chart not found"):
		logger.Warnf("Chart not found in the specified repository")
		return true
	default:
		return false
	}
}
