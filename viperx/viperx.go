package viperx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Description merges the given config file into the Viperx instance (ie. an existing config file)
//
// Notes:
// - It merge the new one to the existing one (which can be even a merge)
// - If the file exists, it is merged into the viper instance
// - If the file does not exist, it does nothing (it keeps the existing viper instance)
// - the first merge => merge the file with nothing
func mergeIfExists(viperx *Viperx, path string, logger logx.Logger) error {
	// merge
	if _, err := os.Stat(path); err == nil {
		viperx.SetConfigFile(path)
		if err := viperx.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return fmt.Errorf("merging config %q: %w", path, err)
			}
		}
	}
	// log when a config file is found
	if _, err := os.Stat(path); err == nil {
		logger.Debugf("found config file: %s", path)
	}
	return nil
}

// Description: returns a Viperx instance scoped to a specific section of the config file
//
// Parameters:
// - fileName: the name of the config file
// - prefix: prefix of the root section in the config file
// - cmdPathName: the rel path of the cobra command
func GetViperx(filename, prefix, cmdPathName string, logger logx.Logger) (*Viperx, error) {

	// define cmd name
	cmdName := filepath.Base(cmdPathName)

	// load the file into a viperx instance
	viperx, err := getViperx(filename, cmdPathName, logger)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// make section items available via sub.xxx
	sub := viperx.Sub(prefix + "." + cmdName)
	if sub == nil {
		return nil, fmt.Errorf("section %q not found", cmdName)
	}
	return &Viperx{Viper: sub}, nil
}

// Description: binds all Cobra command flags to Viper keys so that flags, env vars, and config files work together.
//
// Usage Example:
//  - export GOLUC_WKF_KINDN_EXAMPLE_KEY="env_value"
//  - goluc wkf kindn --example_key="flag_value"

func BindFlags(cmd *cobra.Command, c *Viperx, workflowName string) {
	envPrefix := "GOLUC_WKF"
	c.SetEnvPrefix(envPrefix)
	c.AutomaticEnv()
	c.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		key := f.Name
		if workflowName != "" {
			key = workflowName + "." + key
		}
		if err := c.BindPFlag(key, f); err != nil {
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
