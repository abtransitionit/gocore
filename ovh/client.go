package ovh

import (
	"sync"

	"github.com/abtransitionit/gocore/apicli"
	"github.com/abtransitionit/gocore/logx"
)

// create cached client that support concurent use
var (
	onceOvhClient   sync.Once
	OvhClientCached *apicli.Client

	onceOvhClientToken   sync.Once
	OvhClientTokenCached *apicli.Client
)

// Description: create a client
func GetOvhClientCached(logger logx.Logger) *apicli.Client {
	onceOvhClient.Do(func() {
		OvhClientCached = apicli.NewClient(DOMAIN_EU, logger).
			WithBearerToken(GetAccessTokenFromFileCached)
	})
	return OvhClientCached
}

// Description: create a client
func GetOvhClientTokenCached(logger logx.Logger) *apicli.Client {
	onceOvhClientToken.Do(func() {
		OvhClientTokenCached = apicli.NewClient(DOMAIN_STD, logger)
	})
	return OvhClientTokenCached
}
