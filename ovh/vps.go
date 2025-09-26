package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func GetVpsImageId(vpsNameId string) (string, error) {
	// load VPS list
	listVps, err := GetListVpsFromFile()
	if err != nil {
		return "", fmt.Errorf("failed to load VPS list: %w", err)
	}

	// build a lookup map for NameId -> Distro
	nameIdToDistro := make(map[string]string, len(*listVps))
	for _, vps := range *listVps {
		nameIdToDistro[vps.NameId] = vps.Distro
	}

	// lookup distro by vpsNameId
	distro, ok := nameIdToDistro[vpsNameId]
	if !ok {
		return "", fmt.Errorf("no VPS found with NameId %s", vpsNameId)
	}

	// load images (static map)
	listImage := GetListOsImageFromStruct()

	// lookup image by distro
	image, ok := listImage[distro]
	if !ok {
		return "", fmt.Errorf("no image found for distro %s", distro)
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

	// // inject a dynamic field
	// for key, vps := range listVpsStruct {
	// 	if len(vps.Distro) > 0 {
	// 		vps.NameDynamic = key + strings.ToLower(string(vps.Distro[0]))
	// 		listVpsStruct[key] = vps // reassign updated struct
	// 	}
	// }

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
//   - jsonx.Json:
//   - error
func VpsReinstallHelper(ctx context.Context, logger logx.Logger, vpsNameId string) (jsonx.Json, error) {
	// get ssh key id
	sshKeyId, err := SshKeyGetIdFromFileCached()
	if err != nil {
		return nil, fmt.Errorf("failed to get SSH key id: %w", err)
	}

	// api get ssh key detail
	sshKeyDetail, err := SshKeyGetDetail(ctx, logger, sshKeyId)
	if err != nil {
		return nil, fmt.Errorf("failed to get SSH key detail: %w", err)
	}

	// api get ssh public key
	sshPubKey, err := SshKeyGetPublic(ctx, logger, sshKeyDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to get SSH public key: %w", err)
	}

	// get OS image id
	imageId, err := GetVpsImageId(vpsNameId)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve image id for VPS %s: %w", vpsNameId, err)
	}

	// define the reinstall parameter
	reinstallParam := VpsReinstallParam{
		DoNotSendPassword: true,
		ImageId:           imageId,
		PublicSshKey:      sshPubKey, // example
	}

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
func GetFilteredVpsDetail(ctx context.Context, logger logx.Logger, vpsID, field string) (jsonx.Json, error) {
	// 1 - api get VPS detail
	vpsDetail, err := VpsGetDetail(ctx, logger, vpsID)
	if err != nil {
		return nil, err
	}

	// 2 - Apply optional field filtering directly
	if field != "" {
		val, ok := jsonx.GetField(vpsDetail, field)
		if !ok {
			return nil, fmt.Errorf("field %s not found in VPS detail", field)
		}
		// wrap in jsonx.Json to keep consistent type
		vpsDetail = jsonx.Json{field: val}
	}

	return vpsDetail, nil
}
