package ovh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/abtransitionit/gocore/apicli"
	"github.com/abtransitionit/gocore/filex"
	"github.com/abtransitionit/gocore/logx"
)

// ovh_client.go (OVH-specific)
var cachedToken string
var tokenOnce sync.Once

// read token only if needed
func GetAccessTokenFromFileCached() (string, error) {
	var err error
	tokenOnce.Do(func() {
		cachedToken, err = GetAccessTokenFromFile() // OVH-specific
	})
	return cachedToken, err
}

// Name: CreateAccessTokenForServiceAccount
//
// Description: create an AccessToken from the credential struct.
func CreateAccessTokenForServiceAccount(ctx context.Context, logger logx.Logger) (string, error) {
	// Get service account credentials
	SaId, err := GetSaId()
	if err != nil {
		return "", fmt.Errorf("failed to get client ID: %w", err)
	}
	SaSecret, err := GetSaSecret()
	if err != nil {
		return "", fmt.Errorf("failed to get client secret: %w", err)
	}
	if SaId == "" || SaSecret == "" {
		return "", errors.New("client ID or client secret is not set")
	}

	// create a client without bearer
	client := apicli.NewClient(DOMAIN_STD, logger)

	// define the api action
	ep := endpointReference["CreateToken"]
	endpoint, err := ep.BuildPath(nil)
	if err != nil {
		return "", fmt.Errorf("failed to build path for %s: %w", ep.Desc, err)
	}

	// define the request parameters
	tokenParam := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {SaId},
		"client_secret": {SaSecret},
		"scope":         {"all"},
	}
	reqHeader := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded", // tells the server how the body is encoded
	}
	req := &apicli.Request{
		Verb:     ep.Verb,
		Endpoint: endpoint,
		Headers:  reqHeader,
		Body:     tokenParam,
	}

	// Play the request and get response
	var resp AccessToken
	logger.Infof("%s using endpoint %s", ep.Desc, endpoint)
	err = client.Do(ctx, req, &resp)
	if err != nil {
		return "", fmt.Errorf("API request failed to %s : %w", ep.Desc, err)
	}
	return resp.Token, nil

}

// Name: getCredentialFilePath
func getCredentialFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory %w", err)
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

// Name: getCredentialStrut
//
// Description: get data from a file into  structure
func getCredentialStrut() (*CredentialStruct, error) {
	// define dest structure
	var credStruct CredentialStruct

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

	// Map filecontent into the GO:struct (aka. Unmarshal the JSON)
	if err := json.Unmarshal(fileContent, &credStruct); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// success - return credential as a pointer to a GO struct
	return &credStruct, nil
}

// Name: GetAccessTokenFromFile
//
// Description: gets the AccessToken from the credential struct.
func GetAccessTokenFromFile() (string, error) {
	// get credential as a GO struct
	creds, err := getCredentialStrut()
	if err != nil {
		return "", err
	}
	// success
	return creds.ServiceAccount.AccessToken, nil
}

// Name: CheckTokenExist
//
// Description: checks if the token exists in the credential file.
//
// Returns:
// - an error if the credential file is missing or token is empty.
func CheckTokenExist(ctx context.Context, logger logx.Logger) error {
	logger.Info("Checking token existence")

	// Get the token from the file
	token, err := GetAccessTokenFromFile()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	if token == "" {
		return fmt.Errorf("token is missing or empty")
	}

	logger.Info("Token exists")
	return nil
}

// Name: RefreshToken
//
// Description: generates a new token using the service account and updates the credential file.
//
// Returns
// - the new token or an error if something fails.
func RefreshToken(ctx context.Context, logger logx.Logger) (string, error) {
	logger.Info("Refreshing token")

	// api get a new token from service account
	newToken, err := CreateAccessTokenForServiceAccount(ctx, logger)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	// check
	if newToken == "" {
		return "", fmt.Errorf("obtained token is empty")
	}

	// Save token to file
	if err := updateToken(newToken); err != nil {
		logger.Warnf("could not update token file: %v", err)
	}

	// success
	logger.Info("Token refreshed successfully into credential file")
	return newToken, nil
}

// Name: updateToken
//
// Description: updates the access token in the credential file.
func updateToken(newToken string) error {

	// get file path
	filePath, err := getCredentialFilePath()
	if err != nil {
		return err
	}

	// 1️⃣ Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// 2️⃣ Unmarshal into generic map
	var credentialStruct map[string]interface{}
	// var credentialStruct CredentialStruct
	if err := json.Unmarshal(data, &credentialStruct); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 3️⃣ Update only the token in the JSON file
	sa, ok := credentialStruct["serviceAccount"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("serviceAccount field missing or invalid")
	}
	sa["access_token"] = newToken

	// 4️⃣ Marshal back to JSON
	newData, err := json.MarshalIndent(credentialStruct, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// 5️⃣ Write to file
	if err := os.WriteFile(filePath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
