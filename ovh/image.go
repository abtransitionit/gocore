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

// Description: get the list of images available for a VPS
//
// Notes:
//   - each VPS has its own set of images.
func VpsImageGetList(ctx context.Context, vpsNameOrId string, logger logx.Logger) ([]ImageDetail, error) {
	// 1 - define response structure
	var respData []string

	// 2 - normalize input (can be oxy or vps-xxx)
	vpsId, err := GetVpsId(vpsNameOrId, logger)
	if err != nil {
		return nil, err
	}

	// 3 - define the api endpoint
	action := "VpsImageGetList"
	ep, ok := endpointReference[action]
	if !ok {
		return nil, fmt.Errorf("looking up. Action %q not found", action)
	}

	// 31 - define the endpoint
	endpoint, err := ep.BuildPath(map[string]string{"id": vpsId})
	if err != nil {
		return nil, fmt.Errorf("building path for %s: %w", ep.Desc, err)
	}

	// 4 - define the request
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
	}

	// 5 - create a client
	client := GetOvhClientCached(logger)

	// 6 - Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(ctx, req, &respData)
	if err != nil {
		return nil, fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
	}

	// handle case
	if len(respData) == 0 {
		return nil, fmt.Errorf("no image found")
	}

	// 7 - get details for each image - do it in parallel ang display the result afterwards
	var imageDetailList []ImageDetail
	var wg sync.WaitGroup
	for _, imdItem := range respData {
		wg.Add(1)
		go func(onItem string) {
			defer wg.Done()
			logger.Infof("getting detail for img %s", onItem)
			detail, err := imageGetDetail(ctx, vpsId, onItem, logger)
			if err != nil {
				return
			}
			imageDetailList = append(imageDetailList, *detail)
		}(imdItem)
	} // end for
	wg.Wait()

	// handle success
	return imageDetailList, nil
}

func imageGetDetail(ctx context.Context, vpsNameOrId string, imgId string, logger logx.Logger) (*ImageDetail, error) {
	// 1 - define response structure
	var respData ImageDetail

	// 2 - normalize input (can be oxy or vps-xxx)
	vpsId, err := GetVpsId(vpsNameOrId, logger)
	if err != nil {
		return nil, err
	}

	// 3 - define the api action
	ep := endpointReference["ImageGetDetail"]
	endpoint, err := ep.BuildPath(map[string]string{"idv": vpsId, "idi": imgId})
	if err != nil {
		return nil, fmt.Errorf("building path for %s: %w", ep.Desc, err)
	}

	// 4 - define the request
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
	}

	// 5 - create a client
	client := GetOvhClientCached(logger)

	// 6 - Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(ctx, req, &respData)
	if err != nil {
		return nil, fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
	}

	// success
	return &respData, nil
}
