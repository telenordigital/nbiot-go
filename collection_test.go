package horde

import (
	"net/http"
	"testing"
)

func TestCollection(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	collections, err := client.Collections()
	if err != nil {
		t.Fatal(err)
	}
	if len(collections) != 1 {
		t.Fatalf("expected one collection, got %#v", collections)
	}

	collection, err := client.CreateCollection(Collection{})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteCollection(collection.CollectionID)

	tagKey := "test key"
	tagValue := "test value"
	collection.Tags = map[string]string{tagKey: tagValue}
	collection, err = client.UpdateCollection(collection)
	if err != nil {
		t.Fatal(err)
	}
	if len(collection.Tags) != 1 || collection.Tags[tagKey] != tagValue {
		t.Fatal("unexpected tags:", collection.Tags)
	}

	collections, err = client.Collections()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, c := range collections {
		if c.CollectionID == collection.CollectionID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("collection %v not found in %v", collection, collections)
	}

	if _, err := client.Collection(collection.CollectionID); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteCollection(collection.CollectionID); err != nil {
		t.Fatal(err)
	}
	err = client.DeleteCollection(collection.CollectionID)
	if cerr, ok := err.(ClientError); !ok || cerr.HTTPStatusCode != http.StatusNotFound {
		t.Fatal(err)
	}
}
