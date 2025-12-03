
# Terminology
## domain
- represents **the API server that your client will target**, without the protocol or the endpoint path.
- Example : if the (API) server is accessible via `https://api.example.com/v1/users`, the **domain** is : `api example.com`
- is used to build the client base **URL**. Eg. `urlBase := fmt.Sprintf("https://%s", domain)`
- So all subsequent calls made with this `Client` will use this domain as the base.


# Purpose

- The `apicli` package provides a reusable HTTP client abstraction for interacting with REST APIs. 
- It wraps **Resty** for HTTP requests and supports:
- The package is designed for **progressive development**, allowing you to start with simple GET requests and gradually build more advanced API interactions.

## Features

* Dynamic base URL configuration (`domain`).
* Automatic JSON encoding/decoding.
* Optional bearer token authorization.
* Flexible request construction (headers, query parameters, body).
* Context-aware request execution.
* Endpoint path templating.


## Prerequisites to use the package

Before using this package, ensure you have:

1. **Go environment setup** (Go 1.20+ recommended).

   ```bash
   go version
   ```

2. **Dependencies** are installed:

   * `resty` for HTTP requests:

     ```bash
     go get github.com/go-resty/resty/v2
     ```
   * `logx` (your custom logger interface):

     ```go
     type Logger interface {
         Infof(format string, args ...any)
         Debugf(format string, args ...any)
         Errorf(format string, args ...any)
     }
     ```
   * Standard library packages: `context`, `fmt`, `io`, `net/url`, `strings`.

3. **Understanding**:

   * Basic Go structs, interfaces, and methods.
   * REST API concepts (GET, POST, headers, query parameters, body payloads).


## Initialize a client

**Purpose:** Create a client instance bound to a specific API domain.

```go
import (
    "context"
    "fmt"
    "apicli"
)
// Create a client for a specific API domain
logger := logx.logger 
client := apicli.NewClient("api.example.com", logger)
```

**Explanation:**

* `api.example.com` → base URL for all requests.
* `logger` → optional logging for request lifecycle.
* Default HTTP client (`resty`) is pre-configured with JSON headers.
* Default `context.Background()` is set.


## Add a Bearer Token

**Purpose:** Dynamically inject Authorization token into requests.

```go
client.WithBearerToken(func() (string, error) {
    // Example: fetch token from cache or external service
    return "my-access-token", nil
})
```

* **TokenFunc** is a function returning a token and optional error.
* Token is automatically added to the `Authorization` header before each request.


## Define an Endpoint

Endpoints are represented as `Endpoint` structs:


```go
// define a templated endpoints 
vpsEndpoint := apicli.Endpoint{
    Verb: "GET",
    Path: "/vps/{id}",
    Desc: "Retrieve VPS information by ID",
}
```

**Dynamic Path Parameters:**

```go
// resolved the endpoint with a specific path with parameters
path, err := vpsEndpoint.BuildPath(map[string]string{"id": "12345"})
if err != nil {
    log.Fatal(err)
}
// path = "/vps/12345"
```

* `BuildPath` replaces placeholders like `{id}` with actual values.
* Returns an error if any parameter is empty.


## Create and Execute a Request

**Step 1: Define the request**

```go
# create a request
req := apicli.Request{
    Verb:        "GET",
    Endpoint:    path,
    QueryParams: map[string]string{"page": "1"}, // if any
    Logger:      logger,
}
```

**Step 2: Execute the request**

```go
//Define a response container
var result map[string]interface{}
// Execute the request
err := client.Do(context.Background(), &req, &result)
if err != nil {
    logger.Errorf("Request failed: %v", err)
}
// Print result
fmt.Printf("VPS Info: %+v\n", result)
```

**Explanation:**

* `Do()` executes the HTTP request with context.
* Handles:

  * Context validation.
  * Header injection.
  * Query parameters.
  * Request body (string, `io.Reader`, `url.Values`, or structs/maps for JSON).
* Automatically unmarshals JSON into `out`.

**Error Handling:**

* HTTP 401 → unauthorized token.
* Other 4xx/5xx → returns status and body.
* Successful 2xx → `out` populated with response.



# Progressive Path to create a specific API client

## Step 1: Simple GET

* Initialize a client.
* Make a GET request to a public endpoint.
* Print results.

## Step 2: Path Parameters

* Define endpoints with placeholders (`/vps/{id}`).
* Use `BuildPath()` to populate parameters.

## Step 3: Query Parameters

* Add pagination, filters, or search params via `QueryParams`.

## Step 4: POST Requests

* Send JSON payloads using a struct.
* Observe automatic JSON encoding in `Do()`.

## Step 5: Authorization

* Implement `TokenFunc` to dynamically add bearer tokens.
* Handle token refresh logic.

## Step 6: Logging

* Attach per-request loggers.
* Enable debug-level logs for request/response tracking.

## Step 7: Advanced Features

* Streaming payloads with `io.Reader`.
* Form-encoded requests with `url.Values`.
* Integrate error handling based on HTTP status codes.


## 8. Tips & Best Practices

* Always pass a non-nil `context` to `Do()` to support timeouts and cancellation.
* Reuse `Client` for multiple requests; it maintains configuration.
* Keep endpoints in a `MapEndpoint` for organized management.
* Use `logger.Debugf` for development; switch to `Infof` or `Errorf` in production.
* Wrap `TokenFunc` with caching to avoid repeated token fetching.



# Todo

* Use `context.Context` to support timeouts/cancellation.
* If `body != nil` → marshal to JSON and set `Content-Type: application/json`.
* Always set `Accept: application/json`.
* Merge `c.Headers` into the request (so you can set auth once).
* If `out != nil`, decode `response.Body` into `out` with `json.Decoder`.
* Return useful errors for non-2xx status codes (include status and maybe response snippet).
* Handle `204 No Content` gracefully (don't try to decode).

