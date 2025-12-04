package ovh

import (
	"context"
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/apicli"
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
func VpsGetDetail(ctx context.Context, id string, logger logx.Logger) (jsonx.Json, error) {
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
func vpsReinstall(ctx context.Context, logger logx.Logger, id string, vpsInstallParam VpsReinstallParam) (jsonx.Json, error) {
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
		return nil, fmt.Errorf("API request failed to %s > %w", ep.Desc, err)
	}
	return resp, nil
}

// Description: api reinstall the same OS image on a VPS
//
// Parameters:
//   - ctx: context.Context
//   - logger: logx.Logger
//   - vpsNameOrId: string
//
// Returns:
//   - jsonx.Json: info concerning the reinstalled VPS
//   - error
//
// Notes:
//   - vpsNameOrId can be vps-XXX or o1u, o2a, ...
//   - the returned Json is sended immediately after the API call is received and does not wait the VPS to be ready
func VpsReinstall(ctx context.Context, logger logx.Logger, vpsNameOrId string) (jsonx.Json, error) {
	// 1 - normalize input (can be o1u or vps-xxx)
	vpsId, err := GetVpsId(vpsNameOrId, logger)
	if err != nil {
		return nil, err
	}
	// // 2 - get the Vps:SshKey:id
	// keyId, err := GetSshKeyIdFromFileCached()
	// if err != nil {
	// 	logger.Errorf("failed to get ssh key id %v", err)
	// 	os.Exit(1)
	// }

	// // 1 - get the Vps:SshKey:publicKey
	// sshPubKey, err := SshKeyGetPublicKeyCached(context.Background(), logger, keyId)
	// if err != nil {
	// 	logger.Errorf("getting ssh public key %v", err)
	// 	os.Exit(1)
	// }

	// 3 - get the Vps:OS:ImageId
	imageId, err := GetVpsImageId2(ctx, vpsId, logger)
	if err != nil {
		return nil, fmt.Errorf("resolving image id for VPS %s: %w", vpsId, err)
	}

	// define the reinstall parameter
	reinstallParam := VpsReinstallParam{
		DoNotSendPassword: true,
		ImageId:           imageId,
		PublicSshKey:      sshPubKey, // example
	}

	jsonx.PrettyPrintColor(reinstallParam)

	// reinstall the vps via api
	vpsInfo, err := vpsReinstall(ctx, logger, vpsId, reinstallParam)
	if err != nil {
		return nil, fmt.Errorf("reinstalling VPS %s > %w", vpsId, err)
	}
	return vpsInfo, nil
}

func CheckVpsIsReady(ctx context.Context, logger logx.Logger, vpsId string) (bool, error) {
	// api get vps detail
	vpsDetail, err := VpsGetDetail(ctx, vpsId, logger)
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

// Description: get the VPS id from a name or an id
//
// Notes:
//   - if input is already an Id (e.g., vps-xxx), return it directly
//   - otherwise, treat input as nameDynamic (eg. o1u) and return the corresponding vps-xxx Id in the VPS list
func GetVpsId(vpsNameOrId string, logger logx.Logger) (string, error) {

	// 1 — If input is already an Id (e.g., vps-xxxxxxx.vps.ovh.net), return it directly
	if strings.HasPrefix(vpsNameOrId, "vps-") {
		return vpsNameOrId, nil
	}

	// 2 — Otherwise, treat input as nameDynamic and resolve it
	vpsList, err := GetListVps()
	if err != nil {
		return "", fmt.Errorf("getting list vps: %w", err)
	}

	for _, vps := range *vpsList {
		if vps.NameDynamic == vpsNameOrId {
			return vps.NameId, nil
		}
	}

	return "", fmt.Errorf("getting vps name %q: not found", vpsNameOrId)
}

// func getVpsId(vpsName string, logger logx.Logger) (string, error) {

// 	// 1 - get VPS:list
// 	vpsList, err := GetListVps()
// 	if err != nil {
// 		return "", fmt.Errorf("getting list vps: %w", err)
// 	}

// 	// 2 - iterate through all VPS entries
// 	for _, vps := range *vpsList {
// 		if vps.NameDynamic == vpsName {
// 			return vps.NameId, nil
// 		}
// 	}

// 	return "", fmt.Errorf("getting vps name %q : not found", vpsName)
// }

// // Input can be o1u, o1, nameId, etc.
// func GetVpsId(vpsNameOrId string, logger logx.Logger) (string, error) {
// 	// If input is already an Id (e.g., vps-xxxxxxx.vps.ovh.net)
// 	if strings.HasPrefix(vpsNameOrId, "vps-") {
// 		return vpsNameOrId, nil
// 	}

// 	// Else, treat input as nameDynamic (o1u)
// 	nameId, err := getVpsId(vpsNameOrId, logger)
// 	if err != nil {
// 		return "", err
// 	}

// 	return nameId, nil
// }

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
