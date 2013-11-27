package kinesis

import (
  "net/http"
  "os"
)

type Auth struct {
  AccessKey, SecretKey, Token string
}

// Client is like http.Client, but signs all requests using Auth.
type Client struct {
  Auth *Auth
  // The http client to make requests with. If nil, http.DefaultClient is used.
  Client *http.Client
}

// New creates a new Client.
func NewClient(auth *Auth) *Client {
  if auth.AccessKey == "" {
    auth.AccessKey = os.Getenv("AWS_ACCESS_KEY")
  }
  if auth.SecretKey == "" {
    auth.SecretKey = os.Getenv("AWS_SECRET_KEY")
  }
  return &Client{Auth: auth}
}

func (c *Client) client() *http.Client {
  if c.Client == nil {
    return http.DefaultClient
  }
  return c.Client
}

func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
  Sign(c.Auth, req)
  return c.client().Do(req)
}
