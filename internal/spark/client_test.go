package spark

import "testing"

func TestFetchStatsReturnsErrorWhenEndpointUnavailable(t *testing.T) {
	client := NewClient("http://127.0.0.1:1")
	_, err := client.FetchStats()
	if err == nil {
		t.Fatal("expected an error when the Spark API is unavailable")
	}
}
