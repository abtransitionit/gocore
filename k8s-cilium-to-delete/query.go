package cilium

import (
	"strings"

	"github.com/abtransitionit/gocore/logx"
)

// func HandleCiliumError(err error, logger logx.Logger) bool {
// 	if err == nil {
// 		return false
// 	}
// 	return true
// }

func HandleCiliumError(err error, logger logx.Logger) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "daemonsets.apps") && strings.Contains(msg, "not found"):
		logger.Warnf("Cilium not yet installed in the cluster (daemonset missing)")
		return true
	case strings.Contains(msg, "configmaps") && strings.Contains(msg, "not found"):
		logger.Warnf("Cilium configmap not yet present")
		return true
	default:
		// Don’t swallow true failures
		logger.Debugf("Unhandled Cilium error: %v", err)
		return false
	}
}
