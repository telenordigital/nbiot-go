package nbiot

import (
	"testing"
	"time"
)

func TestCollectionOutputStream(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	collection, err := client.CreateCollection(Collection{})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteCollection(collection.ID)

	stream, err := client.CollectionOutputStream(collection.ID)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			_, err := stream.Recv()
			if err != nil {
				t.Fatalf("%#v", err)
			}
		}
	}()

	time.Sleep(time.Second)
}
