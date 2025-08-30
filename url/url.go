package url

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/abtransitionit/gocore/run"
)

func DownloadLocalArtifact(ctx context.Context, url string, prefix string) (string, error) {
	return DownloadArtifact(ctx, "", url, prefix)
}
func DownloadArtifact(ctx context.Context, vmName, url string, prefix string) (string, error) {
	// Remote download
	if vmName != "" {
		cmd := fmt.Sprintf("goluc do download %s -p %s", url, prefix)
		fileName, err := run.RunCliSsh(vmName, cmd)
		if err != nil {
			return "", fmt.Errorf("failed to remote download file at URL: '%s' from '%s': %w", url, vmName, err)
		}
		return strings.TrimSpace(fileName), nil
	}

	// Local download
	// request the object pointed by the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to GET %s: %w", url, err)
	}

	// clean before exit
	defer resp.Body.Close()

	// Check request status - Only 200 OK is considered successful - Any other code (404, 500, etc.) triggers an error.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status %d when fetching %s", resp.StatusCode, url)
	}

	// Create a temporary file - "gocli-*" → default to a random chars to make it uniq
	tmpFile, err := os.CreateTemp("/tmp", prefix+"-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	// clean before exit
	defer tmpFile.Close()

	// Stream response into the temp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	return tmpFile.Name(), nil
}

// Name: download
//
// Description: fetches a file denoted by an URL AND stores it in a temporary file.
//
// # Parameters
//
// - url: the URL of the file to download
//
// Returns
// - the full path to the temp file or an error if something goes wrong.
func Download(cliName string, url string) (string, error) {

	// request the object pointed by the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to GET %s: %w", url, err)
	}

	// clean before exit
	defer resp.Body.Close()

	// Check request status - Only 200 OK is considered successful - Any other code (404, 500, etc.) triggers an error.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status %d when fetching %s", resp.StatusCode, url)
	}

	// Create a temporary file - "gocli-*" → default to a random chars to make it uniq
	tmpFile, err := os.CreateTemp("/tmp", cliName+"-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	// clean before exit
	defer tmpFile.Close()

	// Stream response into the temp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	return tmpFile.Name(), nil
}
