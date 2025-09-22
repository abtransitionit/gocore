package ovh

import (
	"github.com/abtransitionit/gocore/logx"
	"github.com/ovh/go-ovh/ovh"
)

// Name : OVHClient
//
// Description : wraps the go-ovh client
type OvhClient struct {
	client *ovh.Client
	logger logx.Logger
}

type CredentialStruct struct {
	ServiceAccount struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		AccessToken  string `json:"access_token"`
	} `json:"serviceAccount"`
}

// func main() {
// 	// Path to your credential file
// 	file := os.ExpandEnv("$HOME/wkspc/.config/ovh/credential")

// 	// Read file
// 	data, err := os.ReadFile(file)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Parse JSON
// 	var cfg Config
// 	if err := json.Unmarshal(data, &cfg); err != nil {
// 		panic(err)
// 	}

// 	// Print access_token
// 	fmt.Println(cfg.ServiceAccount.AccessToken)
// }
