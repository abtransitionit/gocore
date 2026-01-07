package ovh

import (
	"context"
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
)

// Description: reads in memory YAML and returns the DistroYaml struct
func GetDistroList() (*DistroYaml, error) {
	// read embeded (aka. cached) file
	theYaml, err := filex.LoadYamlIntoStruct[DistroYaml](yamlDistroList)
	if err != nil {
		return nil, fmt.Errorf("getting YAML file in package > %w", err)
	}
	// return a pointer to the struct
	return theYaml, nil
}

// Description: looks up the DistroName from the DistroId
func GetDistroName(distroId string) (distroName string, err error) {

	// 1 - get list of distro - cached read
	distroList, err := GetDistroList()
	if err != nil {
		return "", err
	}

	// 2 - lookup the field we are searching => return the property: Name
	for _, distro := range distroList.Distro {
		if distro.Id == strings.TrimSpace(distroId) {
			return distro.Name, nil
		}
	}
	return "", fmt.Errorf("not found distro:name for distro:id:%s", distroId) // not found
}

// Description: return the distro Cid from vps
func GetDistroCid(vpsName string) (distroCid string, err error) {

	// 1 - get the list of VPS - cached read
	vpsList, err := GetVpsList()
	if err != nil {
		return "", fmt.Errorf("getting vps:list from configuration file: %v", err)
	}

	// 2 - lookup the field we are searching => return the property: Distro
	for _, item := range vpsList.Vps {
		if item.NameDynamic == strings.TrimSpace(vpsName) {
			return item.Distro, nil
		}
	}
	return "", fmt.Errorf("not found distro:cid for vps:name:%s", vpsName) // not found
}

// Description: return the list of OVH VPS distro available for the vps (api call)
func GetImageList() (*DistroYaml, error) {
	// read embeded (aka. cached) file
	theYaml, err := filex.LoadYamlIntoStruct[DistroYaml](yamlDistroList)
	if err != nil {
		return nil, fmt.Errorf("getting YAML file in package > %w", err)
	}
	// return a pointer to the struct
	return theYaml, nil
}

func GetImageId(ctx context.Context, vpsName string, logger logx.Logger) (string, error) {

	// 1 - normalize input
	vpsId, err := GetVpsId(vpsName, logger)
	if err != nil {
		return "", fmt.Errorf("getting vps id from id %s > %w", vpsName, err)
	}
	// 1 - get distro:cid from vps
	vpsDistroCid, err := GetDistroCid(vpsName)
	if err != nil {
		return "", fmt.Errorf("getting distro cid for id %q > %v", vpsDistroCid, err)
	}

	// 2 - get distro:name from distro:cid
	distroName, err := GetDistroName(vpsDistroCid)
	if err != nil {
		return "", fmt.Errorf("getting distro name from id %q >  %v", vpsDistroCid, err)
	}

	// 3 - get the list of available image for the VPS
	imageDetailList, err := ImageAvailableGetList(ctx, vpsId, logger)
	if err != nil {
		return "", fmt.Errorf("getting vps:distro:list for vps %s > %v", vpsId, err)
	}

	// 4 - lookup the image:id in the list from the distro:name
	for _, item := range imageDetailList {
		if item.Name == distroName {
			// imageId = item.Id
			return item.Id, nil
		}
	}
	return "", nil // not found
	// return imageId, nil
}
