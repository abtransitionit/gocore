package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/abtransitionit/gocore/apicli"
	"github.com/abtransitionit/gocore/filex"
	"github.com/abtransitionit/gocore/jsonx"
	"github.com/abtransitionit/gocore/logx"
)

func VpsGetList(ctx context.Context, logger logx.Logger) ([]string, error) {
	// define response type
	var resp []string

	// define the api action
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
	client := GetOvhClientCached(logger)

	// Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(ctx, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
	}

	// success
	return resp, nil
}
func VpsGetDetail(ctx context.Context, logger logx.Logger, id string) (jsonx.Json, error) {
	// define response type
	var resp jsonx.Json

	// define the api action
	ep := endpointReference["VpsGetDetail"]
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
	if err := client.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
	}
	return resp, nil
}
func VpsGetOs(ctx context.Context, logger logx.Logger, id string) (jsonx.Json, error) {
	// define response type
	var resp jsonx.Json

	// define the api action
	ep := endpointReference["VpsGetOs"]
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

func VpsReinstall(ctx context.Context, logger logx.Logger, id string, vpsInstallParam VpsReinstallParam) (jsonx.Json, error) {
	// define response type
	var resp jsonx.Json

	// define the api action
	ep := endpointReference["VpsReinstall"]
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
		Body:     vpsInstallParam,
	}

	// Play the request and get response
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	if err := client.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
	}
	return resp, nil
}

func GetListOsImageFromStruct() MapVpsOsImage {
	return vpsOsImageReference
}

func GetVpsImageId(vpsNameId string, logger logx.Logger) (string, error) {
	// 1 - get VPS:list (static file)
	listVps, err := GetListVpsFromFileCached()
	if err != nil {
		return "", fmt.Errorf("failed to load VPS list: %w", err)
	}
	// 2 - get OsIage:list (static map)
	listImage := GetListOsImageFromStruct()

	// 3 - lookup the vps which has the vpsNameId
	var vps Vps
	var vpsKey string

	found := false
	// loop over all vps object
	for key, candidate := range *listVps {
		if candidate.NameId == vpsNameId {
			vps = candidate
			vpsKey = key
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("VPS not found with NameId: %s", vpsNameId)
	}
	// here we have the vps
	logger.Infof("vpsNameId: %s, distro: %s, VpsShortName: %s", vpsNameId, vps.Distro, vpsKey)

	// 4 - get the distro
	if strings.TrimSpace(vps.Distro) == "" {
		return "", fmt.Errorf("VPS %s has no distro defined", vpsNameId)
	}

	// 5 - lookup the image ID for this distro
	image, ok := listImage[vps.Distro]
	if !ok {
		return "", fmt.Errorf("no image found for distro: %s", vps.Distro)
	}

	return image.Id, nil

}

func GetListVpsFromFile() (*ListVpsStruct, error) {
	// get List as a GO struct
	listVps, err := getlistVpsStruct()
	if err != nil {
		return nil, err
	}
	// inject a dynamic field into the struct
	for key, vps := range *listVps {
		if len(vps.Distro) > 0 {
			vps.NameDynamic = key + strings.ToLower(string(vps.Distro[0]))
			(*listVps)[key] = vps // reassign updated struct
		}
	}
	// success
	return listVps, nil
}

var (
	listVpsOnce sync.Once
	listVpsVal  *ListVpsStruct
	listVpsErr  error
)

func GetListVpsFromFileCached() (*ListVpsStruct, error) {
	listVpsOnce.Do(func() {
		listVpsVal, listVpsErr = GetListVpsFromFile()
	})
	return listVpsVal, listVpsErr
}

// Name: getCredentialStrut
//
// Description: get file:data into  Go:structure
func getlistVpsStruct() (*ListVpsStruct, error) {
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

// Name: VpsReinstallHelper
//
// Description: api rebuild a VPS
//
// Parameters:
//   - ctx: context.Context
//   - logger: logx.Logger
//   - vpsNameId: string
//
// Returns:
//   - jsonx.Json: info concerning the rebuilded VPS
//   - error
//
// Notes:
//   - the returned Json is sended immediately after the API call is received and does not wait the VPS to be ready
func VpsReinstallHelper(ctx context.Context, logger logx.Logger, vpsNameId string) (jsonx.Json, error) {
	// 1 - get the Vps:SshKey:id
	keyId, err := GetSshKeyIdFromFileCached()
	if err != nil {
		logger.Errorf("failed to get ssh key id %v", err)
		os.Exit(1)
	}

	// 1 - get the Vps:SshKey:publicKey
	sshPubKey, err := SshKeyGetPublicKeyCached(context.Background(), logger, keyId)
	if err != nil {
		logger.Errorf("failed to get ssh public key %v", err)
		os.Exit(1)
	}

	// 3 - get the Vps:OS:ImageId
	imageId, err := GetVpsImageId(vpsNameId, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve image id for VPS %s: %w", vpsNameId, err)
	}

	// define the reinstall parameter
	reinstallParam := VpsReinstallParam{
		DoNotSendPassword: true,
		ImageId:           imageId,
		PublicSshKey:      sshPubKey, // example
	}

	// jsonx.PrettyPrintColor(reinstallParam)

	// reinstall the vps via api
	vpsInfo, err := VpsReinstall(ctx, logger, vpsNameId, reinstallParam)
	if err != nil {
		return nil, fmt.Errorf("failed to reinstall VPS %s: %w", vpsNameId, err)
	}
	return vpsInfo, nil
}

func CheckVpsIsReady(ctx context.Context, logger logx.Logger, vpsNameId string) (bool, error) {
	// api get vps detail
	vpsDetail, err := VpsGetDetail(ctx, logger, vpsNameId)
	if err != nil {
		return false, err
	}
	// get the vps:state
	state, ok := vpsDetail["state"].(string)
	if !ok {
		return false, fmt.Errorf("unexpected state format in VPS detail")
	}

	return state == "running", nil
}

// Name: DisplayVpsDetail
//
// Description: gets a VPS:detail or a VPS:detail:field according to field.
// Returns
// - an error instead of exiting, so the caller can handle it.
// func GetFilteredVpsDetail(ctx context.Context, logger logx.Logger, vpsID, field string) (jsonx.Json, error) {
// 	// 1 - api get VPS detail
// 	vpsDetail, err := VpsGetDetail(ctx, logger, vpsID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 2 - Apply optional field filtering directly
// 	if field != "" {
// 		val, ok := jsonx.GetField(vpsDetail, field)
// 		if !ok {
// 			return nil, fmt.Errorf("field %s not found in VPS detail", field)
// 		}
// 		// wrap in jsonx.Json to keep consistent type
// 		vpsDetail = jsonx.Json{field: val}
// 	}

// 	return vpsDetail, nil
// }
