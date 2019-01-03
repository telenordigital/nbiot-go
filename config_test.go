package nbiot

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFirstMatchingPath(t *testing.T) {
	getFirstMatchingPath(".telenor-nbiot")
}

func TestFileDefaultConfig(t *testing.T) {

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	contents := "address=http://example.com\ntoken=sometoken\n"
	tempFile := "telenor-nbiot.testconfig"
	tempPath := filepath.Join(cwd, "telenor-nbiot.testconfig")

	ioutil.WriteFile(tempPath, []byte(contents), 0666)
	defer os.Remove(tempPath)

	// unset the environment first to make sure it won't interfere with the
	// file. It might contain settings that is in use so set it back to the
	// original value afterwards.
	oldAddr := os.Getenv(AddressEnvironmentVariable)
	oldToken := os.Getenv(TokenEnvironmentVariable)
	defer func() {
		os.Setenv(AddressEnvironmentVariable, oldAddr)
		os.Setenv(TokenEnvironmentVariable, oldToken)
	}()

	os.Setenv(AddressEnvironmentVariable, "")
	os.Setenv(TokenEnvironmentVariable, "")

	address, token, err := addressTokenFromConfig(tempFile)
	if err != nil {
		t.Fatal(err)
	}
	if address != "http://example.com" || token != "sometoken" {
		t.Fatalf("Configuration isn't the expected values: %s / %s", address, token)
	}

	contents = "token=foobar\nsome=thing\nother=thing\n\nsome=thing=other\n\n"
	ioutil.WriteFile(getFirstMatchingPath(tempFile), []byte(contents), 0666)
	_, _, err = addressTokenFromConfig(tempFile)
	if err == nil {
		t.Fatal("expected error")
	}

}

func TestEnvironmentConfig(t *testing.T) {
	oldAddr := os.Getenv(AddressEnvironmentVariable)
	oldToken := os.Getenv(TokenEnvironmentVariable)
	defer func() {
		os.Setenv(AddressEnvironmentVariable, oldAddr)
		os.Setenv(TokenEnvironmentVariable, oldToken)
	}()

	os.Setenv(AddressEnvironmentVariable, "something")
	os.Setenv(TokenEnvironmentVariable, "other")

	address, token, err := addressTokenFromConfig(ConfigFile)
	if err != nil {
		t.Fatal(err)
	}

	if address != "something" {
		t.Fatal("Expected environment variable to override config")
	}
	if token != "other" {
		t.Fatal("Expected environment variable to override config")
	}
}
