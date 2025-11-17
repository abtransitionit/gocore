package property

import (
	"fmt"
	"os"
	"os/user"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetPath(_ ...string) (string, error) {
	path := os.Getenv("PATH")
	if path == "" {
		return "", fmt.Errorf("PATH environment variable is not set")
	}
	return path, nil
}

func GetOsArch(_ ...string) (string, error) {
	return runtime.GOARCH, nil // go env GOARCH
}

func GetOsType(_ ...string) (string, error) {
	return runtime.GOOS, nil // go env GOOS
}

func GetCpu(_ ...string) (string, error) {
	output, err := cpu.Info()
	if err != nil {
		return "", fmt.Errorf("getting cpu info > %v", err)
	}
	return fmt.Sprintf("%v", output[0].Cores), nil
}

func GetRam(_ ...string) (string, error) {
	output, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", output.Total/(1024*1024*1024)), nil
}

func GetOsDistro(_ ...string) (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}
	return info.Platform, nil
}

func GetOsFamily(_ ...string) (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", fmt.Errorf("getting os family > %v", err)
	}
	return info.PlatformFamily, nil
}

func GetOsKernelVersion(_ ...string) (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", fmt.Errorf("getting os kernel version > %v", err)
	}
	return info.KernelVersion, nil
}

func GetOsUser(_ ...string) (string, error) {
	output, err := user.Current()
	if err != nil {
		return "", err
	}
	return output.Username, nil
}

func GetOsVersion(_ ...string) (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", fmt.Errorf("getting os version > %v", err)
	}
	return info.PlatformVersion, nil
}

func GetOsInfos(_ ...string) (string, error) {
	family, err := GetOsFamily()
	if err != nil {
		return "", err
	}

	distro, err := GetOsDistro()
	if err != nil {
		return "", err
	}

	version, err := GetOsVersion()
	if err != nil {
		return "", err
	}

	kernel, err := GetOsKernelVersion()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("family: %-6s :: distro: %-10s :: OsVersion: %-6s :: OsKernelVersion: %s", family, distro, version, kernel), nil
}
