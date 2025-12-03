package ovh

import (
	"context"
	"fmt"
	"os"
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

// Name: VpsReinstallHelper
//
// Description: api reinstall the same OS image on a VPS
//
// Parameters:
//   - ctx: context.Context
//   - logger: logx.Logger
//   - vpsNameId: string
//
// Returns:
//   - jsonx.Json: info concerning the reinstalled VPS
//   - error
//
// Notes:
//   - vpsNameOrId can be vps-XXX or o1u, o2a, ...
//   - the returned Json is sended immediately after the API call is received and does not wait the VPS to be ready
func VpsReinstallHelper(ctx context.Context, logger logx.Logger, vpsNameOrId string) (jsonx.Json, error) {
	// 1 - normalize input (can be o1u or nameId)
	vpsNameId, err := resolveVpsNameId(vpsNameOrId, logger)
	if err != nil {
		return nil, err
	}
	// 2 - get the Vps:SshKey:id
	keyId, err := GetSshKeyIdFromFileCached()
	if err != nil {
		logger.Errorf("failed to get ssh key id %v", err)
		os.Exit(1)
	}

	// 1 - get the Vps:SshKey:publicKey
	sshPubKey, err := SshKeyGetPublicKeyCached(context.Background(), logger, keyId)
	if err != nil {
		logger.Errorf("getting ssh public key %v", err)
		os.Exit(1)
	}

	// 3 - get the Vps:OS:ImageId
	imageId, err := GetVpsImageId(vpsNameId, logger)
	if err != nil {
		return nil, fmt.Errorf("resolving image id for VPS %s: %w", vpsNameId, err)
	}

	// define the reinstall parameter
	reinstallParam := VpsReinstallParam{
		DoNotSendPassword: true,
		ImageId:           imageId,
		PublicSshKey:      sshPubKey, // example
	}

	jsonx.PrettyPrintColor(reinstallParam)

	// reinstall the vps via api
	vpsInfo, err := VpsReinstall(ctx, logger, vpsNameId, reinstallParam)
	if err != nil {
		return nil, fmt.Errorf("reinstalling VPS %s: %w", vpsNameId, err)
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

func getVpsNameId(vpsName string, logger logx.Logger) (string, error) {

	// 1 - get VPS:list (static file, cloned + decorated)
	vpsList, err := GetListVps()
	if err != nil {
		return "", fmt.Errorf("getting list vps: %w", err)
	}

	// 2 - iterate through all VPS entries
	for _, vps := range *vpsList {
		if vps.NameDynamic == vpsName {
			return vps.NameId, nil
		}
	}

	return "", fmt.Errorf("getting vps name %q : not found", vpsName)
}

// Input can be o1u, o1, nameId, etc.
func resolveVpsNameId(input string, logger logx.Logger) (string, error) {
	// If input is already a nameId (e.g., vps-xxxxxxx.vps.ovh.net)
	if strings.HasPrefix(input, "vps-") {
		return input, nil
	}

	// Else, treat input as nameDynamic (o1u)
	nameId, err := getVpsNameId(input, logger)
	if err != nil {
		return "", err
	}

	return nameId, nil
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
