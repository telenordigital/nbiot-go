package nbiot

import (
	"net/http"
	"testing"
)

func TestTeam(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	team, err := client.CreateTeam(Team{})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteTeam(team.ID)

	tagKey := "test key"
	tagValue := "test value"
	team.Tags = map[string]string{tagKey: tagValue}
	team, err = client.UpdateTeam(team)
	if err != nil {
		t.Fatal(err)
	}
	if len(team.Tags) != 1 || team.Tags[tagKey] != tagValue {
		t.Fatal("unexpected tags:", team.Tags)
	}

	teams, err := client.Teams()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, t := range teams {
		if t.ID == team.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("team %v not found in %v", team, teams)
	}

	if _, err := client.Team(team.ID); err != nil {
		t.Fatal(err)
	}

	if ivs, err := client.Invites(team.ID); err != nil {
		t.Fatal(err)
	} else if len(ivs) > 0 {
		t.Fatal(ivs)
	}

	iv, err := client.CreateInvite(team.ID)
	if err != nil {
		t.Fatal(err)
	} else if iv.Code == "" {
		t.Fatal(iv)
	}
	defer client.DeleteInvite(team.ID, iv.Code)

	if ivs, err := client.Invites(team.ID); err != nil {
		t.Fatal(err)
	} else if len(ivs) != 1 || ivs[0] != iv {
		t.Fatal(ivs)
	}

	_, err = client.AcceptInvite(iv.Code)
	if cerr, ok := err.(ClientError); !ok || cerr.HTTPStatusCode != http.StatusConflict {
		t.Fatal(err)
	}

	if err := client.DeleteInvite(team.ID, iv.Code); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteTeam(team.ID); err != nil {
		t.Fatal(err)
	}
	err = client.DeleteTeam(team.ID)
	if cerr, ok := err.(ClientError); !ok || cerr.HTTPStatusCode != http.StatusNotFound {
		t.Fatal(err)
	}
}
