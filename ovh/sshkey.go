package ovh

import (
	"context"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/apicli"
	"github.com/abtransitionit/gocore/jsonx"
	"github.com/abtransitionit/gocore/logx"
)

func SshKeyGetList(ctx context.Context, logger logx.Logger) ([]string, error) {
	// define response type
	var resp []string

	// define the api action
	ep := endpointReference["SshKeyGetList"]
	endpoint, err := ep.BuildPath(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build path for %s: %w", ep.Desc, err)
	}

	// create a client
	client := GetOvhClientCached(logger)

	// define the request structure
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
	}

	// Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(ctx, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
	}
	return resp, nil
}
func SshKeyGetDetail(ctx context.Context, logger logx.Logger, id string) (jsonx.Json, error) {
	// define response type
	var resp jsonx.Json

	// define the api action
	ep := endpointReference["SshKeyGetDetail"]
	endpoint, err := ep.BuildPath(map[string]string{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to build path for %s: %w", ep.Desc, err)
	}

	// create a client
	client := GetOvhClientCached(logger)

	// define the request structure
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
	}

	// Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(ctx, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
	}
	return resp, nil
}
func SshKeyGetDetailCached(ctx context.Context, logger logx.Logger, id string) (jsonx.Json, error) {
	// define response type
	var resp jsonx.Json

	// define the api action
	ep := endpointReference["SshKeyGetDetail"]
	endpoint, err := ep.BuildPath(map[string]string{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to build path for %s: %w", ep.Desc, err)
	}

	// create a client
	client := GetOvhClientCached(logger)

	// define the request structure
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
	}

	// Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(ctx, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
	}
	return resp, nil
}
func SshKeyGetPublic(ctx context.Context, logger logx.Logger, sshKeyDetail jsonx.Json) (string, error) {
	// check parameters
	if sshKeyDetail == nil {
		return "", fmt.Errorf("sshKeyDetail is nil")
	}
	return sshKeyDetail["key"].(string), nil
}

func SshKeyGetIdFromFile() (string, error) {
	// get credential as a GO struct
	creds, err := getCredentialStrut()
	if err != nil {
		return "", err
	}
	// success
	return creds.SshKeyId, nil
}

// cache the result
var (
	cachedSshKeyId   string
	cachedSshKeyOnce sync.Once
	cachedSshKeyErr  error
)

func SshKeyGetIdFromFileCached() (string, error) {
	cachedSshKeyOnce.Do(func() {
		cachedSshKeyId, cachedSshKeyErr = SshKeyGetIdFromFile()
	})
	return cachedSshKeyId, cachedSshKeyErr
}
