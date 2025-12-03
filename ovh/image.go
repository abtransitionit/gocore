package ovh

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/apicli"
	"github.com/abtransitionit/gocore/logx"
)

func GetListImageAvailable() string {
	return ""
}

func ImageGetList(ctx context.Context, logger logx.Logger) ([]ImageDetail, error) {
	// define response type
	var resp []string

	// define the api action
	ep := endpointReference["ImageGetList"]
	endpoint, err := ep.BuildPath(map[string]string{"id": "vps-9c33782a.vps.ovh.net"})
	if err != nil {
		return nil, fmt.Errorf("building path for %s: %w", ep.Desc, err)
	}

	// define the request structure
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
	}

	// create a client
	client := GetOvhClientCached(logger)

	// Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(ctx, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
	}

	// handle case
	if len(resp) == 0 {
		return nil, fmt.Errorf("no image found")
	}

	var images []ImageDetail
	for _, v := range resp {
		// logger.Infof("Found image %s", v)
		detail, err := ImageGetDetail(ctx, v, logger)
		if err != nil {
			return nil, fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
		}
		// logger.Infof("image detail:%s", detail)
		images = append(images, *detail)
	}

	// success
	return images, nil
}

func ImageGetDetail(ctx context.Context, idImage string, logger logx.Logger) (*ImageDetail, error) {
	// define response type
	var resp ImageDetail

	// define the api action
	ep := endpointReference["ImageGetDetail"]
	endpoint, err := ep.BuildPath(map[string]string{"idv": "vps-9c33782a.vps.ovh.net", "idi": idImage})
	if err != nil {
		return nil, fmt.Errorf("building path for %s: %w", ep.Desc, err)
	}

	// define the request structure
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
	}

	// create a client
	client := GetOvhClientCached(logger)

	// Play the request and get response
	// logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(ctx, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
	}

	// success
	return &resp, nil
}
