package gocli

import (
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
)

// func Install(logger logx.Logger, cli GoCli, osType string, osArch string, uname string) (string, error) {
// 	// resolve URL
// 	url, err := ResolveURL(logger, cli, osType, osArch, uname)
// 	if err != nil {
// 		return "", err
// 	}
// 	logger.Infof("Cli: %s Url: %s", cli.Name, url)
// 	return url, nil
// }

// Name:resolveURL
//
// Description: resolves the final download URL for a CLI.
//
// Todos:
// - handle "latest" version resolution here.
func ResolveURL(logger logx.Logger, cli GoCli, osType string, osArch string, uname string) (string, error) {

	// lookup the templated URL in package private database
	template, ok := goCliReference[cli.Name]
	if !ok {
		return "", fmt.Errorf("no cli %s found in Go CLI db", cli.Name)
	}
	if template.Url == "" {
		return "", fmt.Errorf("no URL template defined for %s", cli.Name)
	}
	// For now, just use Version directly
	tag := cli.Version
	return substituteUrlPlaceholders(template.Url, cli, tag, osType, osArch, uname), nil
}

// Name: substituteUrlPlaceholders
//
// Description: replaces placeholders in the URL template.
func substituteUrlPlaceholders(templatedUrl string, cli GoCli, tag string, osType string, osArch string, uname string) string {

	replacements := map[string]string{
		"$NAME":  cli.Name,
		"$TAG":   tag,
		"$OS":    osType,
		"$ARCH":  osArch,
		"$UNAME": uname,
	}
	url := templatedUrl
	for k, v := range replacements {
		url = strings.ReplaceAll(url, k, v)
	}
	return url
}
