package ovh

// Define structs for the request body and response
// The request body needs to be `x-www-form-urlencoded`
// so we'll use a url.Values map and convert it to a string.
type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	// TokenType   string `json:"token_type"`
}

type VpsOsImage struct {
	Name string
	Id   string
}
type MapVpsOsImage map[string]VpsOsImage

type VpsInfo struct {
	DisplayName string
	NameId      string
	Distro      string
}
type MapVpsInfo map[string]VpsInfo

type CredentialStruct struct {
	SshKeyId       string `json:"sshKeyId"`
	ServiceAccount struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		AccessToken  string `json:"access_token"`
	} `json:"serviceAccount"`
}

type ListVpsStruct struct {
	Vps []struct {
		DisplayName string `json:"DisplayName"`
		NameId      string `json:"NameId"`
		Distro      string `json:"Distro"`
	} `json:"vps"`
}

type VpsReinstallParam struct {
	DoNotSendPassword bool   `json:"doNotSendPassword"`
	ImageId           string `json:"imageId"`
	PublicSshKey      string `json:"publicSshKey"`
}
