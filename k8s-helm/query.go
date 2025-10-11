package helm

import (
	"strings"

	"github.com/abtransitionit/gocore/logx"
)

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
