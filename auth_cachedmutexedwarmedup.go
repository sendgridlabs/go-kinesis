package kinesis

import (
	"sync"
	"time"
)

// newCachedMutexedWarmedUpAuth wraps another auth object
// with a cache that is thread-safe, and will always attempt
// to fetch credentials when initialised.
// The underlying Auth object will only be called if the time is
// past the last returned expiration time.
func newCachedMutexedWarmedUpAuth(underlying temporaryCredentialGenerator) (Auth, error) {
	rv := &cachedMutexedAuth{
		underlying: underlying,
	}
	_, err := rv.KeyForSigning(time.Now())
	if err != nil {
		return nil, err
	}
	return rv, nil
}

// Auth interface for authentication credentials and information
type temporaryCredentialGenerator interface {
	// KeyForSigning return an access key / secret / token appropriate for signing at time now,
	// which as the name suggests, is usually now.
	// Additionally returns the expriration time of these credentials.
	ExpiringKeyForSigning(now time.Time) (*SigningKey, time.Time, error)
}

type cachedMutexedAuth struct {
	mu         sync.Mutex
	current    *SigningKey
	expiration time.Time
	underlying temporaryCredentialGenerator
}

func (cmuxa *cachedMutexedAuth) KeyForSigning(now time.Time) (*SigningKey, error) {
	cmuxa.mu.Lock()
	defer cmuxa.mu.Unlock()

	if cmuxa.current == nil || !cmuxa.expiration.After(now) {
		newCurrent, newExpiration, err := cmuxa.underlying.ExpiringKeyForSigning(now)
		if err != nil {
			return nil, err
		}
		cmuxa.current = newCurrent
		cmuxa.expiration = newExpiration
	}

	return cmuxa.current, nil
}
