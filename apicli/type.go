package apicli

import (
	"net/http"
)

// Name: Client
//
// Description: holds shared config (base URL, headers, http client).
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Headers    map[string]string // default headers, e.g. Authorization
}
