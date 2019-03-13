package api

import "fmt"

// Client stores authentication and creates requests
type Client struct {
	Auth Authenticator
}

// NewRequest creates a new authenticated request for the specified url
func (c *Client) NewRequest(URL string, args ...interface{}) *Request {
	endpoint := fmt.Sprintf(URL, args...)

	r := &Request{
		endpoint,
		"GET",
		"",
		make(map[string]string),
	}

	c.Auth.Authenticate(r)

	return r
}
