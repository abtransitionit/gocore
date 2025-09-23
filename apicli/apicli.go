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

	fullURL := fmt.Sprintf("https://%s%s", req.Domain, req.Endpoint)
	_, err := r.Execute(req.Verb, fullURL)
	return err
}
