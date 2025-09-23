package ovh

// Define structs for the request body and response
// The request body needs to be `x-www-form-urlencoded`
// so we'll use a url.Values map and convert it to a string.
type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	// TokenType   string `json:"token_type"`
}

type CredentialStruct struct {
	ServiceAccount struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		AccessToken  string `json:"access_token"`
	} `json:"serviceAccount"`
}
