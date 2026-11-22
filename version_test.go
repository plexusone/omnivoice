package omnivoice

import (
	"encoding/json"
	"os"
	"testing"
)

// TestVersionMatchesChangelog keeps Version in step with the latest release
// in CHANGELOG.json.
func TestVersionMatchesChangelog(t *testing.T) {
	data, err := os.ReadFile("CHANGELOG.json")
	if err != nil {
		t.Fatal(err)
	}
	var cl struct {
		Releases []struct {
			Version string `json:"version"`
		} `json:"releases"`
	}
	if err := json.Unmarshal(data, &cl); err != nil {
		t.Fatal(err)
	}
	if len(cl.Releases) == 0 {
		t.Fatal("CHANGELOG.json has no releases")
	}
	if want := cl.Releases[0].Version; "v"+Version != want {
		t.Errorf("Version = %q, but the latest release in CHANGELOG.json is %s", Version, want)
	}
}
