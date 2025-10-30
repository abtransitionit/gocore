package viperx

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// description: loads the YAML config file and returns a Viper instance
//
// Notes:
//
// searches the file at this location and merge them all in the following order
// - 0 order is package+global+local with overwrite
// - 1 - Package config (cmd/workflow/<workflowName>/conf.yaml)
// - 2 - Global config ($GOLUC_CONFIG env if set elese ~/.config/goluc/workflow/conf.yaml or )
// - 3 - Local config (aka. current working dir ./conf.yaml)
func getConfig() (*viper.Viper, error) {

	// define instance
	v := viper.New()

	// 1 - define package yaml config file location
	_, file, _, _ := runtime.Caller(2) // because it is not called directly but through GetConfigSection
	packagePath := filepath.Join(path.Dir(file), "conf.yaml")
	// 11 - merge (initial load)
	if err := mergeIfExists(v, packagePath); err != nil {
		return nil, err
	}

	// 2 - define global yaml config file location
	globalPath := os.Getenv("GOLUC_CONFIG")
	if globalPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("getting home directory: %w", err)
		}
		globalPath = filepath.Join(homeDir, "wkspc", ".config", "goluc", "workflow", "conf.yaml")
	}
	// 21 - merge
	if err := mergeIfExists(v, globalPath); err != nil {
		return nil, err
	}

	// 3 - define current working dir yaml config file location
	localPath := "conf.yaml"
	// 31 - merge
	if err := mergeIfExists(v, localPath); err != nil {
		return nil, err
	}

	// 4 - check viper instance is not empty
	if len(v.AllKeys()) == 0 {
		return nil, fmt.Errorf("no configuration file loaded (checked %q and %q and %q in working dir)", packagePath, globalPath, localPath)
	}

	return v, nil
}

// Description merges the given config file into the viper instance
//
// Notes:
// - It merge the new one to the existing one (which can be even a merge)
// - If the file exists, it is merged into the viper instance
// - If the file does not exist, it does nothing (it keeps the existing viper instance)
func mergeIfExists(v *viper.Viper, path string) error {
	if _, err := os.Stat(path); err == nil {
		v.SetConfigFile(path)
		if err := v.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return fmt.Errorf("merging config %q: %w", path, err)
			}
		}
	}
	return nil
}

// Description: returns a Viper instance scoped to a specific section of the YAML
func GetConfig(name string) (*Config, error) {
	// load all config
	v, err := getConfig()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// make section items available via sub.xxx
	sub := v.Sub("workflow." + name)
	if sub == nil {
		return nil, fmt.Errorf("section %q not found", name)
	}
	return &Config{Viper: sub}, nil
}

// Description: binds all Cobra command flags to Viper keys so that flags, env vars, and config files work together.
//
// Usage Example:
//  - export GOLUC_WKF_KINDN_EXAMPLE_KEY="env_value"
//  - goluc wkf kindn --example_key="flag_value"

func BindFlags(cmd *cobra.Command, v *Config, workflowName string) {
	envPrefix := "GOLUC_WKF"
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		key := f.Name
		if workflowName != "" {
			key = workflowName + "." + key
		}
		if err := v.BindPFlag(key, f); err != nil {
			panic(err)
		}
	})
}

// _, file, _, _ := runtime.Caller(0)
// fmt.Println("yo Package:", path.Dir(file))
// // If globalPath does not exist, fallback to absolute YAML
// if _, err := os.Stat(globalPath); err != nil {
// 	// globalPath = "/Users/max/wkspc/git/goluc/cmd/workflow/kindn/conf.yaml"
// }
// var defaultConf []byte

// func GetSection2(name string, logger logx.Logger) (*viper.Viper, error) {
// 	v, err := LoadConfig(logger)
// 	if err != nil {
// 		return nil, fmt.Errorf("loading config: %w", err)
// 	}

// 	// Create a new Viper instance for the workflow subtree
// 	sub := viper.New()
// 	workflowKey := "workflow." + name

// 	for _, key := range v.AllKeys() {
// 		if len(key) >= len(workflowKey) && key[:len(workflowKey)] == workflowKey {
// 			sub.Set(key[len(workflowKey)+1:], v.Get(key))
// 		}
// 	}

// 	return sub, nil
// }
// log current package
// pc, file, _, _ := runtime.Caller(0)
// pkg := path.Dir(runtime.FuncForPC(pc).Name())
// logger.Infof("Package: %s", pkg)
// logger.Infof("Package: %s", path.Dir(file))

// pc, _, _, _ := runtime.Caller(0)
// fn := runtime.FuncForPC(pc).Name() // e.g. "github.com/.../cmd/workflow/kindn.init"
// pkg := path.Base(path.Dir(fn))
// logger.Infof("Package: %s", pkg)
// log current package
// pc, file, _, _ := runtime.Caller(0)
// pkg := path.Dir(runtime.FuncForPC(pc).Name())
// logger.Infof("Package: %s", pkg)
// logger.Infof("Package: %s", path.Dir(file))

// pc, _, _, _ := runtime.Caller(0)
// fn := runtime.FuncForPC(pc).Name() // e.g. "github.com/.../cmd/workflow/kindn.init"
// pkg := path.Base(path.Dir(fn))
// logger.Infof("Package: %s", pkg)
