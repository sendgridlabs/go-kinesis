package kinesis

import (
	"time"
)

const (
	AWSSecurityTokenHeader = "X-Amz-Security-Token"
)

// Auth interface for authentication credentials and information
type Auth interface {
	// KeyForSigning return an access key / secret / token appropriate for signing at time now,
	// which as the name suggests, is usually now.
	KeyForSigning(now time.Time) (*SigningKey, error)
}

// SigningKey returns a set of data needed for signing
type SigningKey struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}
