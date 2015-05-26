package kinesis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	ACCESS_ENV_KEY = "AWS_ACCESS_KEY"
	SECRET_ENV_KEY = "AWS_SECRET_KEY"

	AWS_METADATA_SERVER = "169.254.169.254"
	AWS_IAM_CREDS_PATH  = "/latest/meta-data/iam/security-credentials"
	AWS_IAM_CREDS_URL   = "http://" + AWS_METADATA_SERVER + AWS_IAM_CREDS_PATH
)

// Auth interface for authentication credentials and information
type Auth interface {
	GetToken() string
	GetExpiration() time.Time
	GetSecretKey() string
	GetAccessKey() string
	HasExpiration() bool
	Renew() error
	Sign(*Service, time.Time) []byte
}

type auth struct {
	// accessKey, secretKey are the standard AWS auth credentials
	accessKey, secretKey, token string

	// expiry indicates the time at which these credentials expire. If this is set
	// to anything other than the zero value, indicates that the credentials are
	// temporary (and probably fetched from an IAM role from the metadata server)
	expiry time.Time
}

func NewAuth(accessKey, secretKey string) Auth {
	return &auth{
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

// NewAuthFromEnv retrieves auth credentials from environment vars
func NewAuthFromEnv() (Auth, error) {
	accessKey := os.Getenv(ACCESS_ENV_KEY)
	secretKey := os.Getenv(SECRET_ENV_KEY)
	if accessKey == "" {
		return nil, fmt.Errorf("Unable to retrieve access key from %s env variable", ACCESS_ENV_KEY)
	}
	if secretKey == "" {
		return nil, fmt.Errorf("Unable to retrieve secret key from %s env variable", SECRET_ENV_KEY)
	}

	return NewAuth(accessKey, secretKey), nil
}

// NewAuthFromMetadata retrieves auth credentials from the metadata
// server. If an IAM role is associated with the instance we are running on, the
// metadata server will expose credentials for that role under a known endpoint.
//
// TODO: specify custom network (connect, read) timeouts, else this will block
// for the default timeout durations.
func NewAuthFromMetadata() (Auth, error) {
	auth := &auth{}
	if err := auth.Renew(); err != nil {
		return nil, err
	}
	return auth, nil
}

// HasExpiration returns true if the expiration time is non-zero and false otherwise
func (a *auth) HasExpiration() bool {
	return !a.expiry.IsZero()
}

// GetExpiration retrieves the current expiration time
func (a *auth) GetExpiration() time.Time {
	return a.expiry
}

// GetToken returns the token
func (a *auth) GetToken() string {
	return a.token
}

// GetSecretKey returns the secret key
func (a *auth) GetSecretKey() string {
	return a.secretKey
}

// GetAccessKey returns the access key
func (a *auth) GetAccessKey() string {
	return a.accessKey
}

// Renew retrieves a new token and mutates it on an instance of the Auth struct
func (a *auth) Renew() error {
	role, err := retrieveIAMRole()
	if err != nil {
		return err
	}

	data, err := retrieveAWSCredentials(role)
	if err != nil {
		return err
	}

	// Ignore the error, it just means we won't be able to refresh the
	// credentials when they expire.
	expiry, _ := time.Parse(time.RFC3339, data["Expiration"])

	a.expiry = expiry
	a.accessKey = data["AccessKeyId"]
	a.secretKey = data["SecretAccessKey"]
	a.token = data["Token"]
	return nil
}

// Sign API request by
// http://docs.amazonwebservices.com/general/latest/gr/signature-version-4.html

func (a *auth) Sign(s *Service, t time.Time) []byte {
	h := ghmac([]byte("AWS4"+a.GetSecretKey()), []byte(t.Format(iSO8601BasicFormatShort)))
	h = ghmac(h, []byte(s.Region))
	h = ghmac(h, []byte(s.Name))
	h = ghmac(h, []byte(AWS4_URL))
	return h
}

func retrieveAWSCredentials(role string) (map[string]string, error) {
	var bodybytes []byte
	// Retrieve the json for this role
	resp, err := http.Get(AWS_IAM_CREDS_URL + "/" + role)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	bodybytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jsondata := make(map[string]string)
	err = json.Unmarshal(bodybytes, &jsondata)
	if err != nil {
		return nil, err
	}

	return jsondata, nil
}

func retrieveIAMRole() (string, error) {
	var bodybytes []byte

	resp, err := http.Get(AWS_IAM_CREDS_URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", err
	}
	defer resp.Body.Close()

	bodybytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// pick the first IAM role
	role := strings.Split(string(bodybytes), "\n")[0]
	if len(role) == 0 {
		return "", errors.New("Unable to retrieve IAM role")
	}

	return role, nil
}
