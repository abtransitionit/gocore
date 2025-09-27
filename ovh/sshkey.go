package ovh

import (
	"context"
	"fmt"
	"strings"
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

// func SshKeyGetPublicKey2(ctx context.Context, logger logx.Logger, sshKeyDetail jsonx.Json) (string, error) {
// 	// check parameters
// 	if sshKeyDetail == nil {
// 		return "", fmt.Errorf("sshKeyDetail is nil")
// 	}
// 	return sshKeyDetail["key"].(string), nil
// }

func GetSshKeyIdFromFile() (string, error) {
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
	cachedSshKeyOnce sync.Once
	cachedSshKeyErr  error
	cachedSshKeyId   string
)

func GetSshKeyIdFromFileCached() (string, error) {
	cachedSshKeyOnce.Do(func() {
		cachedSshKeyId, cachedSshKeyErr = GetSshKeyIdFromFile()
	})
	return cachedSshKeyId, cachedSshKeyErr
}

func SshKeyGetPublicKey(ctx context.Context, logger logx.Logger, sshKeyId string) (string, error) {
	// 1 - check parameters
	if strings.TrimSpace(sshKeyId) == "" {
		return "", fmt.Errorf("sshKeyId not provided")
	}
	// 2 - api get ssh key detail
	sshKeyDetail, err := SshKeyGetDetail(ctx, logger, sshKeyId)
	if err != nil {
		return "", fmt.Errorf("failed to get SSH key detail from key id: %s: %w", sshKeyId, err)
	}
	// 3 - rertieve the public key as a json
	field := "key"
	sshPubKeyJson, err := jsonx.GetFilteredJson(ctx, logger, sshKeyDetail, field)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve field %s from %s : %v", field, sshKeyDetail, err)
	}
	// 4 - extract the value from the kvpair
	val, ok := sshPubKeyJson[field]
	if !ok {
		return "", fmt.Errorf("field %s not found in JSON", field)
	}
	// 5 - check the value is indeed a string
	sshPubKeyString, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("field %s is not a string", field)
	}

	return sshPubKeyString, nil

}

var (
	sshPubKeyOnce sync.Once
	sshPubKeyErr  error
	sshPubKey     string
)

func SshKeyGetPublicKeyCached(ctx context.Context, logger logx.Logger, sshKeyId string) (string, error) {
	sshPubKeyOnce.Do(func() {
		// the closure execute the function only during the program's lifetime once no matter sshKeyId
		// and always return the same couple (sshPubKey, sshPubKeyErr)
		sshPubKey, sshPubKeyErr = SshKeyGetPublicKey(ctx, logger, sshKeyId)
	})
	// This returns the result stored in the GLOBAL variables.
	return sshPubKey, sshPubKeyErr
}
