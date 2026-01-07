package ovh

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	lold "github.com/abtransitionit/gocore/filex"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/mock/filex"
)

// Description: get the file path of the static file containing the list of VPS
func getListVpsFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory %w", err)
	}

	vpsListFilePath := filepath.Join(home, vpsListRelPath)

	ok, err := lold.ExistsFile(vpsListFilePath)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("credential file not found: %s", vpsListFilePath)
	}

	return vpsListFilePath, nil
}

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
//
// Notes:
//   - this function is thread-safe
//   - if it uses embed the file is already cached and this code is not needed
//   - but if in the future we read from another location (e.g., disk), this code will be useful
func getlistVpsCached() (*VpsYaml, error) {
	listVpsOnce.Do(func() {
		listVpsVal, listVpsErr = getlistVpsAsStruct()
	})
	return listVpsVal, listVpsErr
}

// Description: add/inject a dynamic field to an existing struct
// Notes:
//   - NameDynamic = key + first letter of Distro in lowercase
//   - it uses a copy of the struct to avoid modifying the cached version
//   - it returns the modified copy with the added dynamic field
func (listVps VpsYaml) addDynamicField() VpsYaml {
	// 1 - clone the original
	clone := listVps.clone()

	// 2 - add the dynamic field
	for key, vps := range clone.Vps {
		if len(vps.Distro) > 0 {
			vps.NameDynamic = key + strings.ToLower(string(vps.Distro[0]))
			clone.Vps[key] = vps
		}
	}
	return clone
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

// Description: get the list of organization's VPS (cached file read) with dynamic field injected
func GetVpsList() (*VpsYaml, error) {
	// get file cached
	listVps, err := getlistVpsCached()
	if err != nil {
		return nil, err
	}

	// get the list with dynamic field
	listVpsWithDynamicField := listVps.addDynamicField()

	// success
	return &listVpsWithDynamicField, nil
}

// // Description: get the os image id of a VPS
// func GetVpsImageId(vpsNameId string, logger logx.Logger) (string, error) {
// 	// 1 - get VPS:list
// 	listVps, err := GetVpsList()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to load VPS list: %w", err)
// 	}
// 	// 2 - get OsIage:list (static map)
// 	listImage := GetListOsImageFromStruct()

// 	// 3 - lookup the vps which has the vpsNameId
// 	var vps Vps
// 	var vpsKey string

// 	found := false
// 	// loop over all vps object
// 	for key, candidate := range listVps.Vps {
// 		if candidate.Id == vpsNameId {
// 			vps = candidate
// 			vpsKey = key
// 			found = true
// 			break
// 		}
// 	}
// 	if !found {
// 		return "", fmt.Errorf("VPS not found with NameId: %s", vpsNameId)
// 	}
// 	// here we have the vps
// 	logger.Infof("vpsNameId: %s, distro: %s, VpsShortName: %s", vpsNameId, vps.Distro, vpsKey)

// 	// 4 - get the distro
// 	if strings.TrimSpace(vps.Distro) == "" {
// 		return "", fmt.Errorf("VPS %s has no distro defined", vpsNameId)
// 	}

// 	// 5 - lookup the image ID for this distro
// 	image, ok := listImage[vps.Distro]
// 	if !ok {
// 		return "", fmt.Errorf("no image found for distro: %s", vps.Distro)
// 	}

// 	return image.Id, nil

// }

// Description: return the ovh image id of a vps
//
// Parameters:
//   - vpsName: the vps nameDynamic (e.g., o1u, o2a, ...)
func GetVpsImageId2(ctx context.Context, vpsName string, logger logx.Logger) (string, error) {
	// 1 - get the list of available image for this VPS
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
	vpsDistroAvailableList, err := ImageAvailableGetList(ctx, vpsName, logger)
	if err != nil {
		return "", fmt.Errorf("api getting vps list image available for %s: > %w", vpsName, err)
	}
	logger.Infof("vpsDistroAvailableList: %s", vpsDistroAvailableList)

	// handle success
	return "", nil

}

// // Description: get the list of os image from a static structure
// func GetListOsImageFromStruct() MapVpsOsImage {
// 	return vpsOsImageReference
// }

// Description: get the list of VPS names (cached file read)
//
// Notes:
//   - returns o1u, o2a, ...
func GetVpsListName() ([]string, error) {
	// get file cached
	vpsList, err := GetVpsList()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(vpsList.Vps))
	for _, vps := range vpsList.Vps {
		names = append(names, vps.NameDynamic)
	}
	return names, nil
}

// Description: return the Distro Cid (e.g., debian10, ubuntu2004, ...) of a VPS
//
// Parameters:
//   - name: the VPS nameDynamic (e.g., o1u, o2a, ...)
//
// Returns:
//   - cid mzan Custom Id
func GetVpsDistro(name string) (VpsDistro string, err error) {

	// get file cached
	vpsList, err := GetVpsList()
	if err != nil {
		return "", err
	}

	// lookup by nameDynamic
	for _, vps := range vpsList.Vps {
		if vps.NameDynamic == name {
			return vps.Distro, nil
		}
	}
	return "", nil // not found
}

// Description: return a printable string of the VPS list
//
// Notes:
//   - allow to pretty print the list via the function PrettyPrintTable
func GetPrintableVpsList(vpsList *VpsYaml) string {
	var b strings.Builder

	// header
	b.WriteString("tName\tId\tDistro\n")

	// rows
	for _, vps := range vpsList.Vps {
		fmt.Fprintf(&b, "%s\t%s\t%s\n",
			vps.NameDynamic, vps.Id, vps.Distro)
	}

	return b.String()
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
