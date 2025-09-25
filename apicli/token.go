package apicli

// add a token to client if needed
// func (c *Client) WithBearerToken(token string) *Client {
// 	c.resty.SetHeader("Authorization", "Bearer "+token)
// 	return c
// }

type TokenFunc func() (string, error)

// Notes:
// - WithBearerToken’s job is to register how to get a token, not to actually fetch it immediately.
func (c *Client) WithBearerToken(fn TokenFunc) *Client {
	c.tokenFunc = fn
	return c
}
