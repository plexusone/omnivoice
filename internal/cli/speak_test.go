package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/plexusone/omnivoice"
	"github.com/plexusone/omnivoice-core/tts"
)

// fakeTTS records the last request.
type fakeTTS struct {
	text   string
	config tts.SynthesisConfig
}

func (f *fakeTTS) Name() string { return "fake" }

func (f *fakeTTS) Synthesize(_ context.Context, text string, config tts.SynthesisConfig) (*tts.SynthesisResult, error) {
	f.text, f.config = text, config
	return &tts.SynthesisResult{Audio: []byte("audio:" + tts.ApplyPronunciations(text, config.Pronunciations, true)), Format: config.OutputFormat}, nil
}

func (f *fakeTTS) SynthesizeStream(context.Context, string, tts.SynthesisConfig) (<-chan tts.StreamChunk, error) {
	return nil, nil
}

func (f *fakeTTS) ListVoices(context.Context) ([]tts.Voice, error) {
	return []tts.Voice{
		{ID: "v2", Name: "Marco", Language: "en", Metadata: map[string]any{"accent": "italian"}},
		{ID: "v1", Name: "Alice", Language: "en", Metadata: map[string]any{"accent": "british"}},
	}, nil
}

func (f *fakeTTS) GetVoice(context.Context, string) (*tts.Voice, error) {
	return nil, tts.ErrVoiceNotFound
}

func registerFake(t *testing.T) *fakeTTS {
	t.Helper()
	f := &fakeTTS{}
	omnivoice.RegisterTTSProvider("fake", func(omnivoice.ProviderConfig) (tts.Provider, error) { return f, nil }, omnivoice.PriorityThick)
	t.Setenv("FAKE_API_KEY", "test")
	return f
}

func TestSpeak(t *testing.T) {
	f := registerFake(t)
	dir := t.TempDir()
	terms := filepath.Join(dir, "terms")
	if err := os.MkdirAll(terms, 0o750); err != nil {
		t.Fatal(err)
	}
	term := `{"id": "lamborghini", "canonicalForm": "Lamborghini", "ipa": {"en-US": "lamborˈɡiːni"}, "pronunciations": {"en-US": "lambor gheeny"}}`
	if err := os.WriteFile(filepath.Join(terms, "lamborghini.json"), []byte(term), 0o600); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "a.wav")
	o := &speakOptions{provider: "fake", voice: "v2", output: out, locale: "en-US", termsDir: terms,
		ipa: []string{"Miura=/ˈmjura/"}, say: []string{"Lamborghini = lam bor ghini"}, phonemes: "on"}
	if err := runSpeak(context.Background(), o, nil, strings.NewReader("Lamborghini built the Miura.\n")); err != nil {
		t.Fatal(err)
	}
	if f.text != "Lamborghini built the Miura." || f.config.VoiceID != "v2" || f.config.OutputFormat != "wav" {
		t.Errorf("request = %q, %+v", f.text, f.config)
	}
	if f.config.Extensions["elevenlabs.phonemes"] != true {
		t.Errorf("extensions = %v", f.config.Extensions)
	}
	got, err := os.ReadFile(out) //nolint:gosec // test file
	if err != nil {
		t.Fatal(err)
	}
	want := `audio:<phoneme alphabet="ipa" ph="lamborˈɡiːni">Lamborghini</phoneme> built the <phoneme alphabet="ipa" ph="ˈmjura">Miura</phoneme>.`
	if string(got) != want {
		t.Errorf("audio = %s", got)
	}
	ps := f.config.Pronunciations
	if len(ps) != 2 || ps[0].Respelling != "lam bor ghini" || ps[0].IPA != "lamborˈɡiːni" || ps[1].IPA != "ˈmjura" {
		t.Errorf("pronunciations = %+v", ps)
	}
}

func TestSpeakErrors(t *testing.T) {
	registerFake(t)
	base := speakOptions{provider: "fake", voice: "v", output: filepath.Join(t.TempDir(), "a.mp3"), phonemes: "auto"}
	for name, mod := range map[string]func(*speakOptions){
		"bad --ipa":      func(o *speakOptions) { o.ipa = []string{"Miura"} },
		"bad --phonemes": func(o *speakOptions) { o.phonemes = "yes" },
		"missing terms":  func(o *speakOptions) { o.termsDir = filepath.Join(t.TempDir(), "none") },
		"no key":         func(o *speakOptions) { o.provider = "nokey" },
	} {
		o := base
		mod(&o)
		if err := runSpeak(context.Background(), &o, []string{"text"}, nil); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	if err := runSpeak(context.Background(), &base, nil, strings.NewReader("  ")); err == nil {
		t.Error("empty input: no error")
	}
}

func TestVoices(t *testing.T) {
	registerFake(t)
	var out bytes.Buffer
	if err := runVoices(context.Background(), &voicesOptions{provider: "fake", match: "ITALIAN"}, &out); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); !strings.HasPrefix(got, "v2 ") || strings.Contains(got, "Alice") {
		t.Errorf("voices = %q", got)
	}
}

func TestOutputFormat(t *testing.T) {
	for in, want := range map[string]string{"a.mp3": "mp3", "a.WAV": "wav", "a.opus": "opus", "a": "mp3", "a.m4a": "mp3"} {
		if got := outputFormat(in); got != want {
			t.Errorf("outputFormat(%q) = %q", in, got)
		}
	}
}
