package kinesis

import (
  "net/http"
  "os"
)

const (
  ACCESS_ENV_KEY  = "AWS_ACCESS_KEY"
  SECRET_ENV_KEY  = "AWS_SECRET_KEY"
  REGION_ENV_NAME = "AWS_REGION_NAME"
)

// Auth store information about AWS Credentials
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
    auth.AccessKey = os.Getenv(ACCESS_ENV_KEY)
  }
  if auth.SecretKey == "" {
    auth.SecretKey = os.Getenv(SECRET_ENV_KEY)
  }
  return &Client{Auth: auth}
}

func GetRegion(region Region) string {
  if region.Name == "" {
    return os.Getenv(REGION_ENV_NAME)
  }
  return region.Name
}

// get client
func (c *Client) client() *http.Client {
  if c.Client == nil {
    return http.DefaultClient
  }
  return c.Client
}

// do some request, but sign it before sending
func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
  Sign(c.Auth, req)
  return c.client().Do(req)
}
