package nbiot

import (
	"math/rand"
	"net/http"
	"strconv"
	"testing"
)

func TestDownstream(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	collection, err := client.CreateCollection(Collection{})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteCollection(collection.ID)

	var devices []Device
	for i := 0; i < 5; i++ {
		d, err := client.CreateDevice(collection.ID, Device{
			IMSI: strconv.Itoa(rand.Intn(1e15)),
			IMEI: strconv.Itoa(rand.Intn(1e15)),
		})
		if err != nil {
			t.Fatal(err)
		}
		defer client.DeleteDevice(collection.ID, d.ID)
		devices = append(devices, d)
	}

	err = client.Send(collection.ID, devices[0].ID, DownstreamMessage{Port: 1234, Payload: []byte("Hello, device!")})
	if cerr, ok := err.(ClientError); !ok || cerr.HTTPStatusCode != http.StatusConflict {
		t.Fatal(err)
	}

	if testing.Short() {
		return
	}

	res, err := client.Broadcast(collection.ID, DownstreamMessage{Port: 1234, Payload: []byte("Hello, device!")})
	if err != nil || res.Failed != len(devices) {
		t.Fatal(err, res)
	}
}
