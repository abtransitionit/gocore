package apicli

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
)

// Name: Client
//
// Description: holds shared config (base URL, headers, http client).
// type Client struct {
// 	BaseURL    string
// 	HTTPClient *http.Client
// 	Headers    map[string]string // default headers, e.g. Authorization
// }

type Request struct {
	Verb        string            // eg. "POST", "GET"
	Domain      string            // eg. "www.ovh.com"
	Endpoint    string            // eg. "/auth/oauth2/token"
	Headers     map[string]string // ie. curl:--header
	Body        any               // depends on the content-type header: `application/x-www-form-urlencoded`, `application/json`, `text/plain`, `multipart/form-data`
	QueryParams map[string]string // eg. {"page": "1"}
	Context     context.Context
	Logger      logx.Logger
}

// example
// endpoint="/auth/oauth2/token"
// domain="www.ovh.com"
//  ... so that url=="https://${domain}${endpoint}"
