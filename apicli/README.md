
# Todo

* Use `context.Context` to support timeouts/cancellation.
* If `body != nil` → marshal to JSON and set `Content-Type: application/json`.
* Always set `Accept: application/json`.
* Merge `c.Headers` into the request (so you can set auth once).
* If `out != nil`, decode `response.Body` into `out` with `json.Decoder`.
* Return useful errors for non-2xx status codes (include status and maybe response snippet).
* Handle `204 No Content` gracefully (don't try to decode).

