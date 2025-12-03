package ovh

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/url"
)

func ListInfo() (string, error) {

	logger := logx.GetLogger()
	logger.Infof("🔹 Ovh API Endpoint : %s", ovhEndpoint01)
	logger.Info(url.Display("🔹 Docs", ovhApiDoc))
	logger.Info(url.Display("🔹 Ovh Go lib on Github", ovhGithubGoLib))
	logger.Info(url.Display("🔹 Liste ESN", listPrestatire))
	logger.Info(url.Display("🔹 Explore OVH API(s)", ovhEndpoint02))
	logger.Info(url.Display("🔹 IHM API", ovhApiIhm))
	logger.Infof("🔹 Credential file path: %s", credentialFilePath)
	logger.Infof("🔹 APi: create App via gui: %s", createApp)
	return "", nil
}

func InstallVpsImage(ctx context.Context, hostname string, logger logx.Logger) (string, error) {

	logger.Infof("%s > Re-install the same os image on a VPS", hostname)
	return "", nil
}

func UpdateVpsImage() (string, error) {

	logger := logx.GetLogger()
	logger.Infof("🔹 install another OS image on the VPS")
	return "", nil
}
