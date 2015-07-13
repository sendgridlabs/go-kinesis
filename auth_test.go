package kinesis

import (
	"os"
	"testing"
)

func TestAuthInterfaceIsImplemented(t *testing.T) {
	var auth Auth = &AuthCredentials{}
	if auth == nil {
		t.Error("Invalid nil auth credentials value")
	}
}

func TestGetSecretKey(t *testing.T) {
	auth := NewAuth("BAD_ACCESS_KEY", "BAD_SECRET_KEY")

	if auth.GetAccessKey() != "BAD_ACCESS_KEY" {
		t.Error("incorrect value for auth#accessKey")
	}
}

func TestGetAccessKey(t *testing.T) {
	auth := NewAuth("BAD_ACCESS_KEY", "BAD_SECRET_KEY")

	if auth.GetSecretKey() != "BAD_SECRET_KEY" {
		t.Error("incorrect value for auth#secretKey")
	}
}

func TestNewAuthFromEnv(t *testing.T) {
	os.Setenv(AccessEnvKey, "asdf")
	os.Setenv(SecretEnvKey, "asdf")

	auth, _ := NewAuthFromEnv()

	if auth.GetAccessKey() != "asdf" {
		t.Error("Expected AccessKey to be inferred as \"asdf\"")
	}

	if auth.GetSecretKey() != "asdf" {
		t.Error("Expected SecretKey to be inferred as \"asdf\"")
	}

	// Validate that the fallback environment variables will also work
	os.Setenv(AccessEnvKey, "") // Use Unsetenv with go1.4
	os.Setenv(SecretEnvKey, "") // Use Unsetenv with go1.4
}

func TestNewAuthFromEnvWithFallbackVars(t *testing.T) {
	os.Setenv(AccessEnvKeyId, "asdf")
	os.Setenv(SecretEnvAccessKey, "asdf")

	auth, _ := NewAuthFromEnv()

	if auth.GetAccessKey() != "asdf" {
		t.Error("Expected AccessKey to be inferred as \"asdf\"")
	}

	if auth.GetSecretKey() != "asdf" {
		t.Error("Expected SecretKey to be inferred as \"asdf\"")
	}

	os.Setenv(AccessEnvKeyId, "")
	os.Setenv(SecretEnvAccessKey, "")
}
