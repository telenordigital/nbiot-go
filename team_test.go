package horde

import (
	"net/http"
	"testing"
)

func TestTeam(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	teams, err := client.Teams()
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 1 {
		t.Fatalf("expected one team, got %#v", teams)
	}
	// if len(teams[0].Members) != 1 {
	// 	t.Fatalf("expected one team member, got %#v", teams[0].Members)
	// }

	team, err := client.CreateTeam(Team{})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteTeam(team.TeamID)

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

	teams, err = client.Teams()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, t := range teams {
		if t.TeamID == team.TeamID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("team %v not found in %v", team, teams)
	}

	if _, err := client.Team(team.TeamID); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteTeam(team.TeamID); err != nil {
		t.Fatal(err)
	}
	err = client.DeleteTeam(team.TeamID)
	if cerr, ok := err.(ClientError); !ok || cerr.HTTPStatusCode != http.StatusNotFound {
		t.Fatal(err)
	}
}
