package ovh

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/abtransitionit/gocore/filex"
	"github.com/pkg/errors"
)

// Name: getCredentialFilePath
func getCredentialFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve home directory")
	}

	credentialPath := filepath.Join(home, credentialRelPath)

	ok, err := filex.ExistsFile(credentialPath)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("credential file not found: %s", credentialPath)
	}

	return credentialPath, nil

}

func getCredentialStrut() (*CredentialStruct, error) {
	// get file path
	filePath, err := getCredentialFilePath()
	if err != nil {
		return nil, err
	}

	// Read the entire file content.
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal the JSON data into the Go struct - map json to a GO struct
	var credentialConfigFile CredentialStruct
	if err := json.Unmarshal(fileContent, &credentialConfigFile); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// success - return credential as a pointer to a GO struct
	return &credentialConfigFile, nil
}

// Name: GetSaId
//
// Description: gets the ClientID from the credential struct.
func GetSaId() (string, error) {
	// get credential as a GO struct
	creds, err := getCredentialStrut()
	if err != nil {
		return "", err
	}
	// success
	return creds.ServiceAccount.ClientID, nil
}

// Name: GetSaSecret
//
// Description: gets the ClientSecret from the credential struct.
func GetSaSecret() (string, error) {
	// get credential as a GO struct
	creds, err := getCredentialStrut()
	if err != nil {
		return "", err
	}
	// success
	return creds.ServiceAccount.ClientSecret, nil
}
