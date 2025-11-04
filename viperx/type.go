package viperx

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/abtransitionit/gocore/logx"
	"github.com/spf13/viper"
)

// Description: Viperx is a wrapper around viper.Viper that denotes a config file (YAML, JSON, ...)
type Viperx struct {
	*viper.Viper
}

// ---------- CONSTRUCTOR ----------

// description: constructor that return an instance of Viperx
//
// Parameters:
// - fileName: the name of the YAML config file
// - cmdName: the name of the command
//
// Notes:
//
// searches the file at this location and merge them all in the following order
// - loads the YAML config file and returns an instance of Viperx that denote a config file (YAML, JSON, ...)
// - 0 order is package+global+local with overwrite
// - 1 - Package config (cmd/workflow/<workflowName>/conf.yaml)
// - 2 - Global config ($GOLUC_CONFIG env if set elese ~/.config/goluc/workflow/conf.yaml or )
// - 3 - Local config (aka. current working dir ./conf.yaml)
func getViperx(fileName, cmdPathName string, logger logx.Logger) (*Viperx, error) {

	// create an instance of Viperx
	viperx := &Viperx{viper.New()}

	// 1 - define package yaml config file location
	_, file, _, _ := runtime.Caller(2) // because it is not called directly but through GetConfigSection
	packagePath := filepath.Join(path.Dir(file), "..", cmdPathName, fileName)
	// 11 - merge (initial load)
	if err := mergeIfExists(viperx, packagePath, logger); err != nil {
		return nil, err
	}

	// 2 - define global yaml config file location
	globalPath := os.Getenv("GOLUC_CONFIG")
	if globalPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("getting home directory: %w", err)
		}
		globalPath = filepath.Join(homeDir, "wkspc", ".config", "goluc", "workflow", fileName)
	}
	// 21 - merge
	if err := mergeIfExists(viperx, globalPath, logger); err != nil {
		return nil, err
	}

	// 3 - define current working dir yaml config file location
	localPath := fileName
	// 31 - merge
	if err := mergeIfExists(viperx, localPath, logger); err != nil {
		return nil, err
	}

	// 4 - check viper instance is not empty
	if len(viperx.AllKeys()) == 0 {
		return nil, fmt.Errorf("no configuration file loaded (checked %q and %q and %q in working dir)", packagePath, globalPath, localPath)
	}

	return viperx, nil
}
