package ovh

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/apicli"
	"github.com/abtransitionit/gocore/logx"
)

func VpsList(ctx context.Context, logger logx.Logger) ([]string, error) {
	// get token from file
	accessToken, err := GetAccessTokenFromFile()
	if err != nil {
		logger.Errorf("%v", err)
		return nil, err
	}
	logger.Infof("Loaded token. First char are: %s", accessToken[:10])

	// define var
	domain := DOMAIN_EU
	endpoint := fmt.Sprintf("%s%s", NS_VERSION, "/vps")
	urlBase := fmt.Sprintf("https://%s", domain)

	// define the request structure
	req := &apicli.Request{
		Verb:     "GET",
		Domain:   domain,
		Endpoint: endpoint,
		Headers: map[string]string{
			"Accept":        "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		},
		Context: ctx,
		Logger:  logger,
	}

	// Define the response's struct
	var resp []string

	// Create the client
	client := apicli.NewClient(urlBase)

	// Play the request and get response
	err = client.Do(req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list VPS: %w", err)
	}
	return resp, nil

}
