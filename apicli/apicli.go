package apicli

import (
	"context"
	"net/http"

	"github.com/abtransitionit/gocore/logx"
)

// Name: Do
//
// Description: performs a request to `path` (relative to BaseURL).
//
// Notes:
// - ctx for cancel/timeouts
// - method: "GET", "POST", ...
// - path: "/iam/policy" or "/iam/policy/123"
// - body: optional request payload (marshal to JSON if non-nil)
// - out: optional pointer where response JSON will be decoded (struct or map[string]interface{})
// - It delegates actual work to helpers (buildRequest, sendRequest, decodeResponse).
func (client *Client) Do(
	ctx context.Context,
	method, path string,
	body any,
	out any,
	logger logx.Logger,
	queryParams map[string]string,
	headers map[string]string,
) error {
	//log
	logger.Infof("Do called: %s %s", method, path)

	// Build request
	req, err := client.buildRequest(ctx, method, path, body, logger, queryParams, headers)
	if err != nil {
		return err
	}

	// Send request
	resp, err := client.sendRequest(req, logger)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response
	return client.decodeResponse(resp, out, logger)
}

// Name: NewClient
//
// Description: creates a client
//
// Notes:
// - if httpClient is nil we use http.DefaultClient.
func NewClient(baseURL string, httpClient *http.Client, headers map[string]string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
		Headers:    headers,
	}
}

func (c *Client) buildRequest(
	ctx context.Context,
	method, path string,
	body any,
	logger logx.Logger,
	queryParams map[string]string,
	headers map[string]string,
) (*http.Request, error) {
	logger.Infof("[buildRequest] method=%s path=%s body=%#v query=%#v headers=%#v",
		method, path, body, queryParams, headers)
	return &http.Request{}, nil
}

func (c *Client) sendRequest(req *http.Request, logger logx.Logger) (*http.Response, error) {
	logger.Info("[sendRequest] request object received")
	return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
}

func (c *Client) decodeResponse(resp *http.Response, out any, logger logx.Logger) error {
	logger.Infof("[decodeResponse] status=%d out=%#v", resp.StatusCode, out)
	return nil
}
