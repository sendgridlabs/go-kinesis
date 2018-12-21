package kinesis

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// NewAuthWithAssumedRole will call STS in a given region to assume a role
// stsAuth object is used to authenticate to STS to fetch temporary credentials
// for the desired role.
func NewAuthWithAssumedRole(roleArn, sessionName, region string, stsAuth Auth) (Auth, error) {
	return newCachedMutexedWarmedUpAuth(&stsCreds{
		RoleARN:     roleArn,
		SessionName: sessionName,
		Region:      region,
		STSAuth:     stsAuth,
	})
}

type stsCreds struct {
	RoleARN     string
	SessionName string
	Region      string
	STSAuth     Auth
}

func (sts *stsCreds) ExpiringKeyForSigning(now time.Time) (*SigningKey, time.Time, error) {
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://sts.%s.amazonaws.com/?%s", sts.Region, (url.Values{
		"Version":         []string{"2011-06-15"},
		"Action":          []string{"AssumeRole"},
		"RoleSessionName": []string{sts.SessionName},
		"RoleArn":         []string{sts.RoleARN},
	}).Encode()), bytes.NewReader([]byte{}))
	if err != nil {
		return nil, time.Time{}, err
	}

	err = (&Service{
		Name:   "sts",
		Region: sts.Region,
	}).Sign(sts.STSAuth, r)
	if err != nil {
		return nil, time.Time{}, err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, time.Time{}, errors.New("bad status code")
	}

	var wrapper struct {
		AssumeRoleResult struct {
			Credentials struct {
				AccessKeyId     string
				SecretAccessKey string
				SessionToken    string
				Expiration      time.Time
			}
		}
	}
	err = xml.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return nil, time.Time{}, err
	}

	// sanity check at least 1 field
	if wrapper.AssumeRoleResult.Credentials.SecretAccessKey == "" {
		return nil, time.Time{}, errors.New("bad data back")
	}

	return &SigningKey{
		AccessKeyId:     wrapper.AssumeRoleResult.Credentials.AccessKeyId,
		SecretAccessKey: wrapper.AssumeRoleResult.Credentials.SecretAccessKey,
		SessionToken:    wrapper.AssumeRoleResult.Credentials.SessionToken,
	}, wrapper.AssumeRoleResult.Credentials.Expiration, nil
}
