package kinesis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	AWSMetadataServer = "169.254.169.254"
	AWSIAMCredsPath   = "/latest/meta-data/iam/security-credentials"
	AWSIAMCredsURL    = "http://" + AWSMetadataServer + "/" + AWSIAMCredsPath
)

// NewAuthFromMetadata retrieves auth credentials from the metadata
// server. If an IAM role is associated with the instance we are running on, the
// metadata server will expose credentials for that role under a known endpoint.
//
// TODO: specify custom network (connect, read) timeouts, else this will block
// for the default timeout durations.
func NewAuthFromMetadata() (Auth, error) {
	return newCachedMutexedWarmedUpAuth(&metadataCreds{})
}

type metadataCreds struct{}

func (mc *metadataCreds) ExpiringKeyForSigning(now time.Time) (*SigningKey, time.Time, error) {
	role, err := retrieveIAMRole()
	if err != nil {
		return nil, time.Time{}, err
	}

	data, err := retrieveAWSCredentials(role)
	if err != nil {
		return nil, time.Time{}, err
	}

	expiry, err := time.Parse(time.RFC3339, data["Expiration"])
	if err != nil {
		return nil, time.Time{}, err
	}

	return &SigningKey{
		AccessKeyId:     data["AccessKeyId"],
		SecretAccessKey: data["SecretAccessKey"],
		SessionToken:    data["Token"],
	}, expiry, nil
}

func retrieveAWSCredentials(role string) (map[string]string, error) {
	var bodybytes []byte

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	// Retrieve the json for this role
	resp, err := client.Get(fmt.Sprintf("%s/%s", AWSIAMCredsURL, role))
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

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	resp, err := client.Get(AWSIAMCredsURL)
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
