package nbiot

import "testing"

func TestOutput(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	collection, err := client.CreateCollection(Collection{})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.DeleteCollection(collection.CollectionID); err != nil {
			t.Fatal(err)
		}
	}()

	output, err := client.CreateOutput(collection.CollectionID, WebHookOutput{URL: DefaultAddr})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.DeleteOutput(collection.CollectionID, output.GetID()); err != nil {
			t.Fatal(err)
		}
	}()

	if logs, err := client.OutputLogs(collection.CollectionID, output.GetID()); err != nil || len(logs) != 0 {
		t.Fatal(err, logs)
	}

	if stat, err := client.OutputStatus(collection.CollectionID, output.GetID()); err != nil || (stat != OutputStatus{}) {
		t.Fatal(err, stat)
	}
}
