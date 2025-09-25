package apicli

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
	"github.com/go-resty/resty/v2"
)

// Name: Client
//
// Description: wraps resty.Client
//
// Notes:
// - the client knows resty.Client, it's domain, and how to generate a token with it's tokenFunc
type Client struct {
	resty     *resty.Client
	domain    string          // eg. "www.ovh.com", "eu.api.ovh.com"
	tokenFunc TokenFunc       // optional, used to set Authorization dynamically and return the access token for the request
	Ctx       context.Context // default context
	Logger    logx.Logger     // default logger
}

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
	Endpoint    string            // eg. "/auth/oauth2/token"
	Headers     map[string]string // ie. curl:--header
	Body        any               // depends on the content-type header: `application/x-www-form-urlencoded`, `application/json`, `text/plain`, `multipart/form-data`
	QueryParams map[string]string // eg. {"page": "1"}
	Context     context.Context
	Logger      logx.Logger
}

type Endpoint struct {
	Verb string
	Path string // e.g. "/vps" or "/vps/{id}"
	Desc string // e.g. "Get Vps info"
}

type MapEndpoint map[string]Endpoint
