package config

import "testing"

func TestParseChannels(t *testing.T) {
	got := parseChannels(" a , b ,, c ")
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
	if len(parseChannels("")) != 0 {
		t.Fatal("empty input should give empty slice")
	}
}

func TestIsChannelAllowed(t *testing.T) {
	c := &Config{YouTubeAllowedChannels: []string{"x", "y"}}
	if !c.IsChannelAllowed("x") {
		t.Fatal("x should be allowed")
	}
	if c.IsChannelAllowed("z") {
		t.Fatal("z should not be allowed")
	}
}
