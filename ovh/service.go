package ovh

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
