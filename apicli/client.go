package apicli

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/go-resty/resty/v2"
)

// Notes:
// - When creating a client we:
//   - set the base URL
//   - set default headers.
//   - pass the domain
func NewClient(domain string, logger logx.Logger) *Client {
	urlBase := fmt.Sprintf("https://%s", domain)
	logger.Infof("Creating client with base URL: %s", urlBase) // debug log
	r := resty.New().
		SetBaseURL(urlBase).
		SetHeader("Accept", "application/json") // default for all requests
	return &Client{
		resty:  r,
		domain: domain,
		Logger: logger,
		Ctx:    context.Background(),
	}
}
