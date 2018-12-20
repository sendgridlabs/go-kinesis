package kinesis

import (
	"fmt"
	"os"
	"time"
)

const (
	AccessEnvKey        = "AWS_ACCESS_KEY"
	AccessEnvKeyId      = "AWS_ACCESS_KEY_ID"
	SecretEnvKey        = "AWS_SECRET_KEY"
	SecretEnvAccessKey  = "AWS_SECRET_ACCESS_KEY"
	SecurityTokenEnvKey = "AWS_SECURITY_TOKEN"
)

// NewAuth creates return an auth object that uses static
// credentials which do not automatically renew.
func NewAuth(accessKey, secretKey, token string) Auth {
	return &staticAuth{
		staticCreds: &SigningKey{
			AccessKeyId:     accessKey,
			SecretAccessKey: secretKey,
			SessionToken:    token,
		},
	}
}

// NewAuthFromEnv retrieves auth credentials from environment vars
func NewAuthFromEnv() (Auth, error) {
	accessKey := os.Getenv(AccessEnvKey)
	if accessKey == "" {
		accessKey = os.Getenv(AccessEnvKeyId)
	}

	secretKey := os.Getenv(SecretEnvKey)
	if secretKey == "" {
		secretKey = os.Getenv(SecretEnvAccessKey)
	}

	token := os.Getenv(SecurityTokenEnvKey)

	if accessKey == "" && secretKey == "" && token == "" {
		return nil, fmt.Errorf("No access key (%s or %s), secret key (%s or %s), or security token (%s) env variables were set", AccessEnvKey, AccessEnvKeyId, SecretEnvKey, SecretEnvAccessKey, SecurityTokenEnvKey)
	}
	if accessKey == "" {
		return nil, fmt.Errorf("Unable to retrieve access key from %s or %s env variables", AccessEnvKey, AccessEnvKeyId)
	}
	if secretKey == "" {
		return nil, fmt.Errorf("Unable to retrieve secret key from %s or %s env variables", SecretEnvKey, SecretEnvAccessKey)
	}

	return NewAuth(accessKey, secretKey, token), nil
}

type staticAuth struct {
	staticCreds *SigningKey
}

func (sc *staticAuth) KeyForSigning(now time.Time) (*SigningKey, error) {
	return sc.staticCreds, nil
}
