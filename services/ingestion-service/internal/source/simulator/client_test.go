package simulator

import "testing"

func TestSimulator_GetClientFromFixture(t *testing.T) {
	client, err := NewClientFromFixtures(
		"./testdata/",
	)
	if err != nil {
		t.Fatalf("fail to create client: %v", err)
	}

	for _, event := range client.scoreboard.Events {
		t.Logf("simulated event: %s, start time: %s, status: %s, completed: %t",
			event.ID, event.Date, event.Status.Type.State, event.Status.Type.Completed)
	}

	for key := range client.plays {
		t.Logf("competition %s has %d game events", key, len(client.plays[key].Items))
	}
}
