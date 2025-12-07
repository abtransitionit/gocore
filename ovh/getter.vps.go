package ovh

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/abtransitionit/gocore/filex"
	"github.com/abtransitionit/gocore/logx"
)

// // Description: get the file path of the static file containing the list of VPS
// func getListVpsFilePath() (string, error) {
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to resolve home directory %w", err)
// 	}

// 	vpsListFilePath := filepath.Join(home, vpsListRelPath)

// 	ok, err := filex.ExistsFile(vpsListFilePath)
// 	if err != nil {
// 		return "", err
// 	}
// 	if !ok {
// 		return "", fmt.Errorf("credential file not found: %s", vpsListFilePath)
// 	}

// 	return vpsListFilePath, nil
// }

// Description: get the content of the file into a Go structure
func getlistVpsAsStruct() (*VpsYaml, error) {
	theYaml, err := filex.LoadYamlIntoStruct[VpsYaml](yamlVpsList)
	if err != nil {
		return nil, fmt.Errorf("getting YAML config file in package: %w", err)
	}
	return theYaml, nil

	// // 1 - get file path
	// filePath, err := getListVpsFilePath()
	// if err != nil {
	// 	return nil, err
	// }
	// 2 - return a pointer to the struct
	// return filex.LoadJsonFromFile[ListVpsStruct](filePath)
}

var (
	listVpsOnce sync.Once
	listVpsVal  *VpsYaml
	listVpsErr  error
)

// Description: get the content of the file into a Go structure (cached)
func getlistVpsCached() (*VpsYaml, error) {
	listVpsOnce.Do(func() {
		listVpsVal, listVpsErr = getlistVpsAsStruct()
	})
	return listVpsVal, listVpsErr
}

// Description: add/inject a dynamic field into the struct
func (listVps VpsYaml) addField() {
	// inject the dynamic field to the copy
	for key, vps := range listVps.Vps {
		if len(vps.Distro) > 0 {
			vps.NameDynamic = key + strings.ToLower(string(vps.Distro[0]))
			listVps.Vps[key] = vps
		}
	}
}

// Description: return a clone of the struct
func (l VpsYaml) clone() VpsYaml {
	// clone the original
	clone := make(map[string]Vps, len(l.Vps))
	for k, v := range l.Vps {
		clone[k] = v
	}
	return VpsYaml{Vps: clone}
}

// Description: get the list of VPS
func GetListVps() (*VpsYaml, error) {
	// get file cached
	listVps, err := getlistVpsCached()
	if err != nil {
		return nil, err
	}

	// clone the cache
	listVpsClone := listVps.clone()

	// inject the dynamic field
	listVpsClone.addField()

	// success
	return &listVpsClone, nil
}

// Description: get the os image id of a VPS
func GetVpsImageId(vpsNameId string, logger logx.Logger) (string, error) {
	// 1 - get VPS:list
	listVps, err := GetListVps()
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
	for key, candidate := range *&listVps.Vps {
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

// Description: get the os image id of a VPS
func GetVpsImageId2(ctx context.Context, vpsId string, logger logx.Logger) (string, error) {
	// 1 - get the list of OVH VPS of the organization (cached file read)
	vpsList, err := getVpsList()
	if err != nil {
		return "", fmt.Errorf("getting organization's list of VPS > %w", err)
	}
	logger.Infof("vpsList: %s", vpsList)

	// 2 - get the list of OVH VPS distro used by the organization (cached file read)
	vpsDistroList, err := getVpsDistroList()
	if err != nil {
		return "", fmt.Errorf("getting organization's list of distro > %w", err)
	}
	logger.Infof("vpsDistroList: %s", vpsDistroList)

	//	3 - get the list of OVH VPS distro available for the vps
	vpsDistroAvailableList, err := VpsImageGetList(ctx, vpsId, logger)
	if err != nil {
		return "", fmt.Errorf("api getting vps list image available for %s: > %w", vpsId, err)
	}
	logger.Infof("vpsDistroAvailableList: %s", vpsDistroAvailableList)

	// handle success
	return "", nil

}

// Description: get the list of os image from a static structure
func GetListOsImageFromStruct() MapVpsOsImage {
	return vpsOsImageReference
}

// var (
// 	listVpsOnce sync.Once
// 	listVpsVal  *ListVpsStruct
// 	listVpsErr  error
// )

// func GetListVpsFromFileCached() (*ListVpsStruct, error) {
// 	listVpsOnce.Do(func() {
// 		listVpsVal, listVpsErr = GetListVps()
// 	})
// 	return listVpsVal, listVpsErr
// }
