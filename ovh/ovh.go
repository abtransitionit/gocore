package ovh

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/abtransitionit/gocore/filex"
	"github.com/pkg/errors"
)

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

// func getAccessToken(clientID, clientSecret string) (*oauth2.Token, error) {
func GetAccessTokenFromFile() (string, error) {

	// get credential file
	filePath, err := getCredentialFilePath()
	if err != nil {
		return "", err
	}

	// Read the entire file content.
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal the JSON data into the Go struct.
	var credentialConfigFile CredentialConfigFile
	if err := json.Unmarshal(fileContent, &credentialConfigFile); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// success
	return credentialConfigFile.ServiceAccount.AccessToken, nil
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
