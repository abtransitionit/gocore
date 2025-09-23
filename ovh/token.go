package ovh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/abtransitionit/gocore/apicli"
	"github.com/abtransitionit/gocore/logx"
)

// Name: CreateAccessTokenForServiceAccount
//
// Description: gets the AccessToken from the credential struct.
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

	// define the request structure
	domain := "www.ovh.com"
	urlBase := fmt.Sprintf("https://%s", domain)
	req := &apicli.Request{
		Verb:     "POST",
		Domain:   domain,
		Endpoint: "/auth/oauth2/token",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded", // ie. define the way the body is defined - here body is a set of curl:--data
		},
		Body: url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {SaId},
			"client_secret": {SaSecret},
			"scope":         {"all"},
		},
		Context: ctx,
		Logger:  logger,
	}

	// Define the response's struct
	var resp AccessToken

	// Create the client
	client := apicli.NewClient(urlBase)

	// Play the request and get response
	if err := client.Do(req, &resp); err != nil {
		return "", fmt.Errorf("failed to fetch access token: %w", err)
	}

	// success
	return resp.Token, nil
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

// CheckTokenExist checks if the token exists in the credential file.
// Returns an error if the credential file is missing or token is empty.
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

// RefreshToken generates a new token using the service account and updates the credential file.
// Returns the new token or an error if something fails.
func RefreshToken(ctx context.Context, logger logx.Logger) (string, error) {
	logger.Info("Refreshing token")

	// Get a new token from service account
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
