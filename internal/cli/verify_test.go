package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/plexusone/omnivoice-core/stt/verify"
)

const transcriptFixture = `{
  "$schema": "https://omnivoice.dev/schema/transcript-v1.json",
  "version": "1.0",
  "text": "we can have the model run all the tests",
  "duration_ms": 4000,
  "segments": [{"text": "we can have the model run all the tests", "start_ms": 0, "end_ms": 4000,
    "words": [
      {"text": "we", "start_ms": 0, "end_ms": 400}, {"text": "can", "start_ms": 400, "end_ms": 800},
      {"text": "have", "start_ms": 800, "end_ms": 1200}, {"text": "the", "start_ms": 1200, "end_ms": 1600},
      {"text": "model", "start_ms": 1600, "end_ms": 2000}, {"text": "run", "start_ms": 2000, "end_ms": 2400},
      {"text": "all", "start_ms": 2400, "end_ms": 2800}, {"text": "the", "start_ms": 2800, "end_ms": 3200},
      {"text": "tests", "start_ms": 3200, "end_ms": 4000}
    ]}],
  "metadata": {"provider": "deepgram", "model": "nova-3", "created_at": "2026-11-17T00:00:00Z"}
}`

func writeFile(t *testing.T, name, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestVerifyAbsence(t *testing.T) {
	tr := writeFile(t, "t.json", transcriptFixture)
	var out bytes.Buffer
	o := &verifyAbsenceOptions{transcripts: []string{tr}, phrases: []string{"run all the tests", "locally"}, minScore: 0.8}
	if err := runVerifyAbsence(o, &out); err != nil {
		t.Fatal(err)
	}
	var res verify.AbsenceResult
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Phrases[0].FoundIn != 1 || res.Phrases[1].FoundIn != 0 || res.Phrases[0].Searches[0].Engine != "deepgram" {
		t.Errorf("result %+v", res)
	}
}

func TestLoadQuotes(t *testing.T) {
	qf := writeFile(t, "q.json", `[{"id": "a", "text": "run all the tests"}]`)
	o := &verifyQuotesOptions{quotesFile: qf, quotes: []string{"the model"}, index: "index.json"}
	qs, err := o.loadQuotes()
	if err != nil {
		t.Fatal(err)
	}
	if len(qs) != 2 || qs[1].ID != "quote-1" {
		t.Errorf("quotes %+v", qs)
	}

	o = &verifyQuotesOptions{quotes: []string{"x"}}
	if _, err := o.loadQuotes(); err == nil || !strings.Contains(err.Error(), "--index") {
		t.Errorf("want an error asking for --index, got %v", err)
	}
	dup := writeFile(t, "d.json", `[{"id": "a", "text": "x", "start_ms": 1, "end_ms": 2}, {"id": "a", "text": "y", "start_ms": 1, "end_ms": 2}]`)
	o = &verifyQuotesOptions{quotesFile: dup}
	if _, err := o.loadQuotes(); err == nil {
		t.Error("want an error for a repeated id")
	}
}

func TestVerifyEnginesNeedKeysAndPlainEngine(t *testing.T) {
	t.Setenv("DEEPGRAM_API_KEY", "")
	o := &verifyQuotesOptions{providers: "deepgram"}
	if _, err := o.engines(false); err == nil || !strings.Contains(err.Error(), "DEEPGRAM_API_KEY") {
		t.Errorf("want a missing-key error, got %v", err)
	}
	t.Setenv("DEEPGRAM_API_KEY", "test")
	t.Setenv("ELEVENLABS_API_KEY", "test")
	o = &verifyQuotesOptions{providers: "deepgram,elevenlabs", plainEngine: "openai"}
	if _, err := o.engines(true); err == nil {
		t.Error("want an error when the plain engine is not a provider")
	}
	o = &verifyQuotesOptions{providers: "deepgram, elevenlabs", models: []string{"deepgram=nova-3"}}
	es, err := o.engines(true)
	if err != nil {
		t.Fatal(err)
	}
	if es[0].Config.Model != "nova-3" || es[0].NoKeywords || !es[1].NoKeywords {
		t.Errorf("engines %+v", es)
	}
	o = &verifyQuotesOptions{providers: "deepgram", models: []string{"nova-3"}}
	if _, err := o.engines(false); err == nil {
		t.Error("want an error for --model without provider=")
	}
}
