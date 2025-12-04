package ovh

import (
	"context"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/apicli"
	"github.com/abtransitionit/gocore/logx"
)

func GetListImageAvailable() string {
	return ""
}

func VpsImageGetList(ctx context.Context, vpsNameOrId string, logger logx.Logger) ([]ImageDetail, error) {
	// define response type
	var resp []string

	// 1 - normalize input (can be o1u or nameId)
	vpsId, err := GetVpsId(vpsNameOrId, logger)
	if err != nil {
		return nil, err
	}

	// 2 - define the api action
	action := "VpsImageGetList"
	ep, ok := endpointReference[action]
	if !ok {
		return nil, fmt.Errorf("looking up. Action %q not found", action)
	}
	// 3 - define the endpoint
	endpoint, err := ep.BuildPath(map[string]string{"id": vpsId})
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
	var wg sync.WaitGroup
	// nbItem := len(resp)
	for _, v := range resp {
		wg.Add(1)
		go func(v string) {
			defer wg.Done()
			logger.Infof("requesting image %s", v)
			detail, err := imageGetDetail(ctx, v, logger)
			if err != nil {
				return
			}
			images = append(images, *detail)
		}(v)
	} // end for
	wg.Wait()

	// success
	return images, nil
}

func imageGetDetail(ctx context.Context, idImage string, logger logx.Logger) (*ImageDetail, error) {
	// define response type
	var resp ImageDetail

	// define the api action
	ep := endpointReference["ImageGetDetail"]
	endpoint, err := ep.BuildPath(map[string]string{"idv": "vps-a7a8f7f6.vps.ovh.net", "idi": idImage})
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
