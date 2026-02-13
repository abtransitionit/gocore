package helm

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
)

// Description: Login with an existing token using a credential file.
//
// Returns
// - an error if something fails.
func LoginWithToken(ctx context.Context, logger logx.Logger) error {
	logger.Info("login using token")

	// get token from credential file
	// token, err := CreateAccessTokenForServiceAccount(ctx, logger)
	// if err != nil {
	// 	return fmt.Errorf("failed to refresh token: %w", err)
	// }

	// mock
	token := "mock-token"
	// check
	if token == "" {
		return fmt.Errorf("no token found in credential file")
	}

	// login to registry with token
	// echo $lRegistryAccessToken | helm registry login $lRegistryDns -u $lRegistryUser --password-stdin

	if token == "" {
		return fmt.Errorf("no token found in credential file")
	}

	// success
	logger.Info("Login successfully from credential file")
	return nil
}
