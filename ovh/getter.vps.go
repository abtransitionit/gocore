package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/abtransitionit/gocore/filex"
	"github.com/abtransitionit/gocore/logx"
)

// Description: get the file path of the static file containing the list of VPS
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

// Description: get the content of the file into a Go structure
func getlistVpsAsStruct() (*ListVpsStruct, error) {
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

var (
	listVpsOnce sync.Once
	listVpsVal  *ListVpsStruct
	listVpsErr  error
)

// Description: get the content of the file into a Go structure (cached)
func getlistVpsCached() (*ListVpsStruct, error) {
	listVpsOnce.Do(func() {
		listVpsVal, listVpsErr = getlistVpsAsStruct()
	})
	return listVpsVal, listVpsErr
}

// Description: add/inject a dynamic field into the struct
func (listVps ListVpsStruct) addField() {
	// inject the dynamic field to the copy
	for key, vps := range listVps {
		if len(vps.Distro) > 0 {
			vps.NameDynamic = key + strings.ToLower(string(vps.Distro[0]))
			listVps[key] = vps
		}
	}
}

// Description: return a clone of the struct
func (l ListVpsStruct) clone() ListVpsStruct {
	// clone the original
	clone := make(ListVpsStruct)
	for k, v := range l {
		clone[k] = v
	}
	return clone
}

// Description: get the list of VPS
func GetListVps() (*ListVpsStruct, error) {
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

// Description: get the os image id of a VPS
func GetVpsImageId2(ctx context.Context, vpsId string, logger logx.Logger) (string, error) {
	var image VpsOsImage
	// var vpsDetail map[string]any
	// 1 - get VPS:Image:list:Available
	// vpsImageList, err := VpsImageGetList(ctx, vpsId, logger)
	// if err != nil {
	// 	return "", fmt.Errorf("api getting vps list image available for %s: > %w", vpsId, err)
	// }
	// // jsonx.PrettyPrintColor(vpsImageList)

	vpsImageListDb, err := getVpsImageList()
	if err != nil {
		return "", fmt.Errorf("api getting vps list image available for %s: > %w", vpsId, err)
	}
	logger.Infof("vpsImageName: %s", vpsImageListDb)

	// // loop over all images
	// for _, vpsImage := range vpsImageList {
	// 	logger.Infof("vpsImageName: %s", vpsImage.Name)
	// }

	// distro, ok := vpsDetail["state"].(string)
	// if !ok {
	// 	return false, fmt.Errorf("unexpected state format in VPS detail")
	// }

	// handle success
	return image.Id, nil

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
