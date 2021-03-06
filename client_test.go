package nbiot

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if _, err := New(); err != nil {
		fmt.Println("Error creating client:", err)
		fmt.Println("You might have to configure nbiot-go via a configuration file or environment variables")
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	if _, err := New(); err != nil {
		t.Fatal("Unable to create the NB-IoT client. You might have to configure it either through environment variables or through a configuration file")
	}
}

func TestSystemDefaults(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.SystemDefaults()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewWithAddress(t *testing.T) {
	address, token, err := addressTokenFromConfig(ConfigFile)
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewWithAddr(address, token)
	if err != nil {
		t.Fatal(err)
	}
}
