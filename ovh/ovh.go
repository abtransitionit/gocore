package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/abtransitionit/gocore/apicli"
	"github.com/abtransitionit/gocore/filex"
	"github.com/abtransitionit/gocore/logx"
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

// func getAccessToken2(clientID, clientSecret string) (*oauth2.Token, error) {
// 	conf := &clientcredentials.Config{
// 		ClientID:     clientID,
// 		ClientSecret: clientSecret,
// 		TokenURL:     "https://api.ovh.com/oauth2/token",
// 	}

// 	client := conf.Client(context.Background())
// 	resp, err := client.Get("https://api.ovh.com/1.0/me")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	// Handle the response
// 	// ...

// 	return nil, nil
// }

// func testOvhApi() {
// 	// Replace with your OAuth2 bearer token
// 	token := "YOUR_OAUTH2_ACCESS_TOKEN"

// 	// Create a client using OAuth2
// 	client := ovh.NewAccessTokenClient("ovh-eu", token)

// 	// Example: get account details
// 	var vps map[string]interface{}
// 	url := "/vps"
// 	err := client.Get(url, &vps)
// 	if err != nil {
// 		log.Fatalf("Error calling OVH API: %v", err)
// 	}

// 	fmt.Println("Vps info:")
// 	fmt.Println(vps)
// }

// "golang.org/x/oauth2"
// "golang.org/x/oauth2/clientcredentials"

// Name: getCredentialStrut
//
// Description: get the credential datas from file into a GO struct.
//
// Returns:
// - *CredentialStruct: a pointer to the a populated CredentialStruct or an error if anything fails.
// - error: an error if anything fails.
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

// GetSaId gets the ClientID from the credential struct.
func GetSaId() (string, error) {
	// get credential as a GO struct
	creds, err := getCredentialStrut()
	if err != nil {
		return "", err
	}
	// success
	return creds.ServiceAccount.ClientID, nil
}

// GetSaSecret gets the ClientSecret from the credential struct.
func GetSaSecret() (string, error) {
	// get credential as a GO struct
	creds, err := getCredentialStrut()
	if err != nil {
		return "", err
	}
	// success
	return creds.ServiceAccount.ClientSecret, nil
}

// GetAccessTokenFromFile gets the AccessToken from the credential struct.
func GetAccessTokenFromFile() (string, error) {
	// get credential as a GO struct
	creds, err := getCredentialStrut()
	if err != nil {
		return "", err
	}
	// success
	return creds.ServiceAccount.AccessToken, nil
}

func GetAccessTokenFromServiceAccount(ctx context.Context, logger logx.Logger) (string, error) {
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

	// Play the request
	if err := client.Do(req, &resp); err != nil {
		return "", fmt.Errorf("failed to fetch access token: %w", err)
	}

	// Save token to file
	if err := updateToken(resp.Token); err != nil {
		logger.Warnf("could not update token file: %v", err)
		// don't return error here, token retrieval was still successful
	}

	return resp.Token, nil
}

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
	var cred map[string]interface{}
	if err := json.Unmarshal(data, &cred); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 3️⃣ Update the token
	sa, ok := cred["serviceAccount"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("serviceAccount field missing or invalid")
	}
	sa["access_token"] = newToken

	// 4️⃣ Marshal back to JSON
	newData, err := json.MarshalIndent(cred, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// 5️⃣ Write file
	if err := os.WriteFile(filePath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// // Create the client
// apiClient := apicli.NewClient("https://www.ovh.com") // base URL optional; domain can also be provided in Request

// // Prepare the request
// req := &apicli.Request{
// 	Verb:     http.MethodPost,
// 	Domain:   "www.ovh.com",
// 	Endpoint: "/auth/oauth2/token",
// 	Body:     body,
// 	Context:  ctx,
// 	Headers: map[string]string{
// 		"Content-Type": "application/x-www-form-urlencoded", // optional, Resty auto-handles this
// 	},
// 	Logger: logger,
// }

// // Response struct
// var resp AccessToken

// // Execute the request
// if err := apiClient.Do(req, &resp); err != nil {
// 	return "", fmt.Errorf("failed to get access token: %w", err)
// }

// return resp.AccessToken, nil
