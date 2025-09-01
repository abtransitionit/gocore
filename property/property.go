// file: gocore/property/property.go
package property

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

var coreProperties = map[string]PropertyHandler{
	"ostype":    getOsType, // e.g. linux, windows, darwin
	"osarch":    getOsArch, // e.g. arm64, amd64
	"cpu":       getCpu,
	"path":      getPath,
	"osversion": getOsVersion,
	"osuser":    getOsUser,
	"ram":       getRam,
	// "netip":      getNetIp,      // code change from original
	// "netgateway": getNetGateway, // code change from original
	"oskversion": getOsKernelVersion,
	"envar":      getEnvar,
}

// GetCorePropertyMap exposes the map of cross-platform properties to external callers.
func GetCorePropertyMap() map[string]PropertyHandler {
	return coreProperties
}

// getOsType retrieves the value of any environment variable
func getEnvar(params ...string) (string, error) {
	if len(params) != 1 {
		return "", fmt.Errorf("expected exactly 1 environment variable name, got %d", len(params))
	}

	// get input
	envarName := params[0]

	return os.Getenv(envarName), nil
}

// getOsType retrieves the operating system type.
func getOsType(_ ...string) (string, error) {
	return strings.ToLower(runtime.GOOS), nil
}

// getOsArch retrieves the operating system architecture.
func getOsArch(_ ...string) (string, error) {
	return runtime.GOARCH, nil
}

// getCpu retrieves the number of CPU cores.
func getCpu(_ ...string) (string, error) {
	output, err := cpu.Info()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", output[0].Cores), nil
}

// getPath retrieves the system PATH environment variable.
func getPath(_ ...string) (string, error) {
	path := os.Getenv("PATH")
	if path == "" {
		return "", fmt.Errorf("PATH environment variable is not set")
	}
	return path, nil
}

// getOsVersion retrieves the platform version.
func getOsVersion(_ ...string) (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}
	return info.PlatformVersion, nil
}

// getOsUser retrieves the current username.
func getOsUser(_ ...string) (string, error) {
	output, err := user.Current()
	if err != nil {
		return "", err
	}
	return output.Username, nil
}

// getRam retrieves the total system RAM in GB.
func getRam(_ ...string) (string, error) {
	output, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", output.Total/(1024*1024*1024)), nil
}

// getNetIp retrieves the system's public IP address.
// func getNetIp(_ ...string) (string, error) {
// 	addrs, err := net.Addrs(true)
// 	if err != nil {
// 		return "", err
// 	}
// 	if len(addrs) > 0 {
// 		return addrs[0].Addr, nil
// 	}
// 	return "unknown", nil
// }

// getNetGateway retrieves the default network gateway.
// func getNetGateway(_ ...string) (string, error) {
// 	routes, err := net.Routes(true)
// 	if err != nil {
// 		return "", err
// 	}
// 	if len(routes) > 0 {
// 		return routes[0].Gateway, nil
// 	}
// 	return "unknown", nil
// }

// getOsKernelVersion retrieves the OS kernel version.
func getOsKernelVersion(_ ...string) (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}
	return info.KernelVersion, nil
}
