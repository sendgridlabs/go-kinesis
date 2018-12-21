package kinesis

import (
	"os"
	"testing"
	"time"
)

func TestGetSecretKey(t *testing.T) {
	auth := NewAuth("BAD_ACCESS_KEY", "BAD_SECRET_KEY", "BAD_SECURITY_TOKEN")
	sk, _ := auth.KeyForSigning(time.Now())
	if sk.AccessKeyId != "BAD_ACCESS_KEY" {
		t.Error("incorrect value for auth#accessKey")
	}
}

func TestGetAccessKey(t *testing.T) {
	auth := NewAuth("BAD_ACCESS_KEY", "BAD_SECRET_KEY", "BAD_SECURITY_TOKEN")
	sk, _ := auth.KeyForSigning(time.Now())
	if sk.SecretAccessKey != "BAD_SECRET_KEY" {
		t.Error("incorrect value for auth#secretKey")
	}
}

func TestGetToken(t *testing.T) {
	auth := NewAuth("BAD_ACCESS_KEY", "BAD_SECRET_KEY", "BAD_SECURITY_TOKEN")
	sk, _ := auth.KeyForSigning(time.Now())
	if sk.SessionToken != "BAD_SECURITY_TOKEN" {
		t.Error("incorrect value for auth#token")
	}
}

func TestNewAuthFromEnv(t *testing.T) {
	os.Setenv(AccessEnvKey, "asdf")
	os.Setenv(SecretEnvKey, "asdf2")
	os.Setenv(SecurityTokenEnvKey, "dummy_token")
	// Validate that the fallback environment variables will also work
	defer os.Unsetenv(AccessEnvKey)
	defer os.Unsetenv(SecretEnvKey)
	defer os.Unsetenv(SecurityTokenEnvKey)

	auth, _ := NewAuthFromEnv()
	sk, _ := auth.KeyForSigning(time.Now())

	if sk.AccessKeyId != "asdf" {
		t.Error("Expected AccessKey to be inferred as \"asdf\"")
	}

	if sk.SecretAccessKey != "asdf2" {
		t.Error("Expected SecretKey to be inferred as \"asdf2\"")
	}

	if sk.SessionToken != "dummy_token" {
		t.Error("Expected SecurityToken to be inferred as \"dummy_token\"")
	}
}

func TestNewAuthFromEnvWithoutSecurityToken(t *testing.T) {
	os.Setenv(AccessEnvKey, "asdf")
	os.Setenv(SecretEnvKey, "asdf2")
	os.Unsetenv(SecurityTokenEnvKey)
	// Validate that the fallback environment variables will also work
	defer os.Unsetenv(AccessEnvKey)
	defer os.Unsetenv(SecretEnvKey)

	auth, _ := NewAuthFromEnv()
	sk, _ := auth.KeyForSigning(time.Now())

	if sk.AccessKeyId != "asdf" {
		t.Error("Expected AccessKey to be inferred as \"asdf\"")
	}

	if sk.SecretAccessKey != "asdf2" {
		t.Error("Expected SecretKey to be inferred as \"asdf2\"")
	}

	if sk.SessionToken != "" {
		t.Error("Expected SecurityToken to be an empty string")
	}
}

func TestNewAuthFromEnvWithoutVars(t *testing.T) {
	os.Unsetenv(AccessEnvKey)
	os.Unsetenv(SecretEnvKey)
	os.Unsetenv(SecurityTokenEnvKey)

	auth, err := NewAuthFromEnv()

	if auth != nil {
		t.Error("Expected auth instance to be nil but was non-nil")
	}

	if err == nil {
		t.Error("Expected error to be non-nil but was nil")
	}
}
func TestNewAuthFromEnvWithFallbackVars(t *testing.T) {
	os.Setenv(AccessEnvKeyId, "asdf")
	os.Setenv(SecretEnvAccessKey, "asdf2")
	os.Setenv(SecurityTokenEnvKey, "dummy_token")
	defer os.Unsetenv(AccessEnvKey)
	defer os.Unsetenv(SecretEnvKey)
	defer os.Unsetenv(SecurityTokenEnvKey)

	auth, _ := NewAuthFromEnv()
	sk, _ := auth.KeyForSigning(time.Now())

	if sk.AccessKeyId != "asdf" {
		t.Error("Expected AccessKey to be inferred as \"asdf\"")
	}

	if sk.SecretAccessKey != "asdf2" {
		t.Error("Expected SecretKey to be inferred as \"asdf2\"")
	}
}
