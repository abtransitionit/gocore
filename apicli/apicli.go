package apicli

import (
	"fmt"
	"io"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// Client wraps resty.Client
type Client struct {
	resty *resty.Client
}

func NewClient(baseURL string) *Client {
	r := resty.New().SetBaseURL(baseURL)
	return &Client{resty: r}
}

func (c *Client) Do(req *Request, out any) error {
	r := c.resty.R().
		SetContext(req.Context).
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

	// define the url to be played
	fullURL := fmt.Sprintf("https://%s%s", req.Domain, req.Endpoint)

	// Execute the request
	resp, err := r.Execute(req.Verb, fullURL)
	if err != nil {
		// Network, context timeout, or other low-level error
		return fmt.Errorf("request execution failed: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode() >= 400 {
		if resp.StatusCode() == 401 {
			// Unauthorized → token may be expired or invalid
			return fmt.Errorf("unauthorized: token may be expired or invalid (HTTP 401)")
		}
		// Other HTTP errors
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	// At this point, status is 2xx and 'out' has been unmarshaled if JSON
	return nil
}
