package api

import (
	"encoding/base64"
	"fmt"
)

type Authenticator interface {
	// Authenticate performs necessary authentication procedures for r such as adding headers
	Authenticate(r *Request)
}

// BasicAuth is a struct containing a username & password to authenticate with
type BasicAuth struct {
	username string
	password string
}

// NewBasicAuth creates an http BasicAuth authenticator
func NewBasicAuth(u, p string) BasicAuth {
	auth := BasicAuth{
		username: u,
		password: p,
	}

	return auth
}

// Authenticate adds the Authorization header
func (auth BasicAuth) Authenticate(r *Request) {
	combined := fmt.Sprintf("%s:%s", auth.username, auth.password)
	b64 := encodeB64([]byte(combined))

	authHeader := fmt.Sprintf("Basic %s", b64)
	r.SetHeader("Authorization", authHeader)
}

// encodeB64 returns the base64 representation of bytes
func encodeB64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}
