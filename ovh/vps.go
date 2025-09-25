package ovh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/abtransitionit/gocore/apicli"
	"github.com/abtransitionit/gocore/filex"
	"github.com/abtransitionit/gocore/jsonx"
	"github.com/abtransitionit/gocore/logx"
)

func VpsGetList(ctx context.Context, logger logx.Logger) ([]string, error) {
	// define response type
	var resp []string

	// define the action
	ep := endpointReference["VpsGetList"]
	endpoint, err := ep.BuildPath(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build path for %s: %w", ep.Desc, err)
	}

	// define the request structure
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
	}

	// create a client
	client := apicli.NewClient(DOMAIN_EU, logger).WithBearerToken(GetCachedAccessToken)

	// Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to %s %w", ep.Desc, err)
	}
	return resp, nil
}
func VpsGetDetail(ctx context.Context, logger logx.Logger, id string) (jsonx.Json, error) {
	// define response type
	var resp jsonx.Json

	// define the action
	ep := endpointReference["VpsGetDetail"]
	endpoint, err := ep.BuildPath(map[string]string{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to build path for %s: %w", ep.Desc, err)
	}

	// create a client
	client := apicli.NewClient(DOMAIN_EU, logger).WithBearerToken(GetCachedAccessToken)

	// define the request structure
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
	}

	// Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	if err := client.Do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to %s %w", ep.Desc, err)
	}
	return resp, nil
}
func VpsGetOs(ctx context.Context, logger logx.Logger, id string) (jsonx.Json, error) {
	// define response type
	var resp jsonx.Json

	// define the action
	ep := endpointReference["VpsGetOs"]
	endpoint, err := ep.BuildPath(map[string]string{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to build path for %s: %w", ep.Desc, err)
	}

	// create a client
	client := apicli.NewClient(DOMAIN_EU, logger).WithBearerToken(GetCachedAccessToken)

	// define the request structure
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
	}

	// Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to %s %w", ep.Desc, err)
	}
	return resp, nil
}

func VpsReinstall(ctx context.Context, logger logx.Logger, id string, vpsInstallParam VpsReinstallParam) (jsonx.Json, error) {
	// define response type
	var resp jsonx.Json

	// define the action
	ep := endpointReference["VpsReinstall"]
	endpoint, err := ep.BuildPath(map[string]string{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to build path for %s: %w", ep.Desc, err)
	}

	// create a client
	client := apicli.NewClient(DOMAIN_EU, logger).WithBearerToken(GetCachedAccessToken)

	// define the request structure
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
		Body:     vpsInstallParam,
	}

	// Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	if err := client.Do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to %s %w", ep.Desc, err)
	}
	return resp, nil
}

func GetOsImageId(VpsNameId string) (string, error) {
	var vps VpsInfo
	// 1️⃣ find VPS by NameId
	found := false
	for _, info := range vpsReference {
		if info.NameId == VpsNameId {
			vps = info
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("VPS with NameId %q not found", VpsNameId)
	}

	// 2️⃣ lookup OS image by Distro
	image, ok := vpsOsImageReference[vps.Distro]
	if !ok {
		return "", fmt.Errorf("Distro %q not found in vpsOsImageReference", vps.Distro)
	}

	// 3️⃣ check if Id is empty
	if image.Id == "" {
		return "", errors.New("OS image Id is empty")
	}

	return image.Id, nil
}

func GetListVpsFromFile() (string, error) {
	// get List as a GO struct
	creds, err := getCredentialStrut()
	if err != nil {
		return "", err
	}
	// success
	return creds.ServiceAccount.AccessToken, nil
}

// Name: getCredentialStrut
//
// Description: get file:data into  Go:structure
func getlistVpsStrut() (*ListVpsStruct, error) {
	// define dest structure
	var listVpsStruct ListVpsStruct

	// get file path
	filePath, err := getListVpsFilePath()
	if err != nil {
		return nil, err
	}

	// Read the entire file content.
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Map filecontent into the GO:struct (aka. Unmarshal the JSON)
	if err := json.Unmarshal(fileContent, &listVpsStruct); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// success - return credential as a pointer to a GO struct
	return &listVpsStruct, nil
}

func getListVpsFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory %w", err)
	}

	listVpsPath := filepath.Join(home, listVpsRelPath)

	ok, err := filex.ExistsFile(listVpsPath)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("credential file not found: %s", listVpsPath)
	}

	return listVpsPath, nil
}
