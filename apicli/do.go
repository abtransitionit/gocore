package apicli

import (
	"context"
	"fmt"
	"io"
	"net/url"
)

func (c *Client) Do(ctx context.Context, req *Request, out any) error {
	if ctx == nil {
		return fmt.Errorf("context must not be nil")
	}

	// Use request logger or fallback to client default
	logger := req.Logger
	if logger == nil && c.Logger != nil {
		logger = c.Logger
	}
	// Create the request
	r := c.resty.R().
		SetContext(ctx).
		SetHeaders(req.Headers).
		SetQueryParams(req.QueryParams).
		SetResult(out)

	// Handle body depending on type
	switch b := req.Body.(type) {
	case url.Values:
		r.SetFormDataFromValues(b)
	case io.Reader:
		r.SetBody(b) // raw stream
	case string:
		r.SetBody(b) // plain text
	default:
		if b != nil {
			r.SetBody(b) // will JSON-encode struct/map
		}
	}

	// Inject token if token function is provided
	if c.tokenFunc != nil {
		token, err := c.tokenFunc()
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}
		r.SetHeader("Authorization", "Bearer "+token)
	}

	// define the url to be played
	fullURL := fmt.Sprintf("https://%s%s", c.domain, req.Endpoint)

	// Execute the request
	resp, err := r.Execute(req.Verb, fullURL)
	if err != nil {
		return fmt.Errorf("request execution failed: %w", err)
	}

	// Handle HTTP errors
	if resp.StatusCode() >= 400 {
		if resp.StatusCode() == 401 {
			// Unauthorized → token may be expired or invalid
			return fmt.Errorf("unauthorized: token may be expired or invalid (HTTP 401)")
		}
		// Other HTTP errors
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	// At this point, status is 2xx and 'out' has been unmarshaled if JSON
	if logger != nil {
		logger.Debugf("Request successful: %s %s", req.Verb, fullURL)
	}
	return nil
}
