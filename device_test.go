package nbiot

import (
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestDevice(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	collection, err := client.CreateCollection(Collection{})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.DeleteCollection(collection.ID); err != nil {
			t.Fatal(err)
		}
	}()

	devices, err := client.Devices(collection.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) != 0 {
		t.Fatalf("expected zero device, got %#v", devices)
	}

	device, err := client.CreateDevice(collection.ID, Device{
		IMSI: strconv.Itoa(rand.Intn(1e15)),
		IMEI: strconv.Itoa(rand.Intn(1e15)),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteDevice(collection.ID, device.ID)

	tagKey := "test key"
	tagValue := "test value"
	imei := strconv.Itoa(rand.Intn(1e15))
	imsi := strconv.Itoa(rand.Intn(1e15))
	device.Tags = map[string]string{tagKey: tagValue}
	device.IMEI = imei
	device.IMSI = imsi
	device, err = client.UpdateDevice(collection.ID, device)
	if err != nil {
		t.Fatal(err)
	}
	if len(device.Tags) != 1 || device.Tags[tagKey] != tagValue {
		t.Fatal("unexpected tags:", device.Tags)
	}
	if device.IMEI != imei || device.IMSI != imsi {
		t.Fatal("unexpected IMEI or IMSI:", device.IMEI, device.IMSI)
	}

	devices, err = client.Devices(collection.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, d := range devices {
		if d.ID == device.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("device %v not found in %v", device, devices)
	}

	if _, err := client.Device(collection.ID, device.ID); err != nil {
		t.Fatal(err)
	}

	data, err := client.DeviceData(collection.ID, devices[0].ID, time.Time{}, time.Time{}, 0)
	if err != nil || len(data) != 0 {
		t.Fatal(err, data)
	}

	if err := client.DeleteDevice(collection.ID, device.ID); err != nil {
		t.Fatal(err)
	}
	err = client.DeleteDevice(collection.ID, device.ID)
	if cerr, ok := err.(ClientError); !ok || cerr.HTTPStatusCode != http.StatusNotFound {
		t.Fatal(err)
	}
}
