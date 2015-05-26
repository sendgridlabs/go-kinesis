package kinesis

import (
	"net/http"
	"time"
)

const AWS_SECURITY_TOKEN_HEADER = "X-Amz-Security-Token"

// Client is like http.Client, but signs all requests using Auth.
type Client struct {
	// Auth holds the credentials for this client instance
	auth *Auth
	// The http client to make requests with. If nil, http.DefaultClient is used.
	client *http.Client
}

// NewClient creates a new Client that uses the credentials in the specified
// Auth object.
//
// This function assumes the Auth object has been sanely initialized. If you
// wish to infer auth credentials from the environment, refer to NewAuth
func NewClient(auth *Auth) *Client {
	return &Client{auth: auth, client: http.DefaultClient}
}

func NewClientWithHTTPClient(auth *Auth, httpClient *http.Client) *Client {
	return &Client{auth: auth, client: httpClient}
}

// Do some request, but sign it before sending
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	err := Sign(c.auth, req)
	if err != nil {
		return nil, err
	}

	if c.auth.HasExpiration() && time.Now().After(c.auth.GetExpiration()) {
		if err = c.auth.Renew(); err != nil { // TODO: (see auth.go#Renew) may be slow
			return nil, err
		}
	}

	if c.auth.GetToken() != "" {
		req.Header.Add(AWS_SECURITY_TOKEN_HEADER, c.auth.GetToken())
	}

	return c.client.Do(req)
}
