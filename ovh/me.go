package ovh

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/apicli"
	"github.com/abtransitionit/gocore/jsonx"
	"github.com/abtransitionit/gocore/logx"
)

func MeGetInfo(ctx context.Context, logger logx.Logger) (jsonx.Json, error) {
	// create a client
	client := apicli.NewClient(DOMAIN_EU, logger).WithBearerToken(GetCachedAccessToken)

	// define the action
	ep := endpointReference["MeGetInfo"]
	endpoint, err := ep.BuildPath(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build path for %s: %w", ep.Desc, err)
	}

	// define the request structure
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
	}

	// Play the request and get response
	var resp jsonx.Json
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to %s %w", ep.Desc, err)
	}
	return resp, nil

}
