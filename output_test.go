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
		if err := client.DeleteCollection(collection.ID); err != nil {
			t.Fatal(err)
		}
	}()

	output, err := client.CreateOutput(collection.ID, WebHookOutput{URL: DefaultAddr})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.DeleteOutput(collection.ID, output.GetID()); err != nil {
			t.Fatal(err)
		}
	}()

	if logs, err := client.OutputLogs(collection.ID, output.GetID()); err != nil || len(logs) != 0 {
		t.Fatal(err, logs)
	}

	if stat, err := client.OutputStatus(collection.ID, output.GetID()); err != nil || (stat != OutputStatus{}) {
		t.Fatal(err, stat)
	}
}
