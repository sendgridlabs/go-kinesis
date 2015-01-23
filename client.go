package kinesis

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	ACCESS_ENV_KEY  = "AWS_ACCESS_KEY"
	SECRET_ENV_KEY  = "AWS_SECRET_KEY"
	REGION_ENV_NAME = "AWS_REGION_NAME"

	AWS_METADATA_SERVER = "169.254.169.254"
	AWS_IAM_CREDS_PATH  = "/latest/meta-data/iam/security-credentials"
	AWS_IAM_CREDS_URL   = "http://" + AWS_METADATA_SERVER + AWS_IAM_CREDS_PATH

	AWS_SECURITY_TOKEN_HEADER = "X-Amz-Security-Token"
)

// Auth store information about AWS Credentials
type Auth struct {
	// AccessKey, SecretKey are the standard AWS auth credentials
	AccessKey, SecretKey, Token string
	// Expiry indicates the time at which these credentials expire. If this is set
	// to anything other than the zero value, indicates that the credentials are
	// temporary (and probably fetched from an IAM role from the metadata server)
	Expiry time.Time
}

// Client is like http.Client, but signs all requests using Auth.
type Client struct {
	// Auth holds the credentials for this client instance
	Auth *Auth
	// The http client to make requests with. If nil, http.DefaultClient is used.
	Client *http.Client
}

// NewAuth returns a new Auth object whose members (AccessKey, SecretKey, etc)
// have been initialized by inspecting the environment or querying the AWS
// metadata server (in that order).
func NewAuth() (auth Auth) {
	// first try grabbing the credentials from the environment
	if auth.AccessKey == "" || auth.SecretKey == "" {
		auth.InferCredentialsFromEnv()
	}

	// if they're still not set, try the metadata server
	if auth.AccessKey == "" || auth.SecretKey == "" {
		auth.InferCredentialsFromMetadata()
	}

	return
}

// InferCredentialsFromEnv retrieves auth credentials from environment vars
func (auth *Auth) InferCredentialsFromEnv() {
	auth.AccessKey = os.Getenv(ACCESS_ENV_KEY)
	auth.SecretKey = os.Getenv(SECRET_ENV_KEY)
}

// InferCredentialsFromMetadata retrieves auth credentials from the metadata
// server. If an IAM role is associated with the instance we are running on, the
// metadata server will expose credentials for that role under a known endpoint.
//
// TODO: specify custom network (connect, read) timeouts, else this will block
// for the default timeout durations.
func (auth *Auth) InferCredentialsFromMetadata() {
	resp1, err := http.Get(AWS_IAM_CREDS_URL)
	if err != nil || resp1.StatusCode != http.StatusOK {
		return
	}
	defer resp1.Body.Close()

	bodybytes, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return
	}

	// pick the first IAM role
	role := strings.Split(string(bodybytes), "\n")[0]
	if len(role) == 0 {
		return
	}

	// Retrieve the json for this role
	resp2, err := http.Get(AWS_IAM_CREDS_URL + "/" + role)
	if err != nil || resp2.StatusCode != http.StatusOK {
		return
	}
	defer resp2.Body.Close()

	bodybytes, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		return
	}

	jsondata := make(map[string]string)
	err = json.Unmarshal(bodybytes, &jsondata)
	if err != nil {
		return
	}

	expiry, _ := time.Parse(time.RFC3339, jsondata["Expiration"])
	// Ignore the error, it just means we won't be able to refresh the
	// credentials when they expire.

	auth.Expiry = expiry
	auth.AccessKey = jsondata["AccessKeyId"]
	auth.SecretKey = jsondata["SecretAccessKey"]
	auth.Token = jsondata["Token"]
}

// NewClient creates a new Client that uses the credentials in the specified
// Auth object.
//
// This function assumes the Auth object has been sanely initialized. If you
// wish to infer auth credentials from the environment, refer to NewAuth
func NewClient(auth *Auth) *Client {
	return &Client{Auth: auth}
}

// GetRegion returns the region name string
func GetRegion(region Region) string {
	if region.Name == "" {
		return os.Getenv(REGION_ENV_NAME)
	}
	return region.Name
}

// get the http client we use to communicate with the server
func (c *Client) client() *http.Client {
	if c.Client == nil {
		return http.DefaultClient
	}
	return c.Client
}

// Do some request, but sign it before sending
func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	Sign(c.Auth, req)

	if !c.Auth.Expiry.IsZero() {
		if time.Now().After(c.Auth.Expiry) {
			c.Auth.InferCredentialsFromMetadata() // TODO: (see above) may be slow
		}
	}

	if len(c.Auth.Token) != 0 {
		req.Header.Add(AWS_SECURITY_TOKEN_HEADER, c.Auth.Token)
	}

	return c.client().Do(req)
}
