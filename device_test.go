package horde

import (
	"net/http"
	"testing"
)

func TestDevice(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	collection, err := client.CreateCollection(Collection{})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteCollection(collection.CollectionID)

	devices, err := client.Devices(collection.CollectionID)
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) != 0 {
		t.Fatalf("expected zero device, got %#v", devices)
	}

	device, err := client.CreateDevice(collection.CollectionID, Device{
		IMEI: str("12"),
		IMSI: str("34"),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteDevice(collection.CollectionID, *device.DeviceID)

	tagKey := "test key"
	tagValue := "test value"
	imei := "56"
	imsi := "78"
	device.Tags = map[string]string{tagKey: tagValue}
	device.IMEI = &imei
	device.IMSI = &imsi
	device, err = client.UpdateDevice(collection.CollectionID, device)
	if err != nil {
		t.Fatal(err)
	}
	if len(device.Tags) != 1 || device.Tags[tagKey] != tagValue {
		t.Fatal("unexpected tags:", device.Tags)
	}
	if device.IMEI == nil || device.IMSI == nil || *device.IMEI != imei || *device.IMSI != imsi {
		t.Fatal("unexpected IMEI or IMSI:", device.IMEI, device.IMSI)
	}

	devices, err = client.Devices(collection.CollectionID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, d := range devices {
		if *d.DeviceID == *device.DeviceID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("device %v not found in %v", device, devices)
	}

	if _, err := client.Device(collection.CollectionID, *device.DeviceID); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteDevice(collection.CollectionID, *device.DeviceID); err != nil {
		t.Fatal(err)
	}
	err = client.DeleteDevice(collection.CollectionID, *device.DeviceID)
	if cerr, ok := err.(ClientError); !ok || cerr.HTTPStatusCode != http.StatusNotFound {
		t.Fatal(err)
	}
}

func str(s string) *string { return &s }
