package cli

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
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
	audio := []byte("audio:" + tts.ApplyPronunciations(text, config.Pronunciations, true))
	return &tts.SynthesisResult{Audio: audio, Format: config.OutputFormat, SampleRate: 16000}, nil
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

// registerFake registers the fake provider, which applies pronunciations
// itself, as ElevenLabs does.
func registerFake(t *testing.T) *fakeTTS {
	t.Helper()
	f := &fakeTTS{}
	omnivoice.RegisterTTSProvider("fake", func(omnivoice.ProviderConfig) (tts.Provider, error) { return f, nil }, omnivoice.PriorityThick)
	t.Setenv("FAKE_API_KEY", "test")
	appliesPronunciations["fake"] = true
	t.Cleanup(func() { delete(appliesPronunciations, "fake") })
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
	if len(got) < 44 || string(got[:4]) != "RIFF" || string(got[44:]) != want {
		t.Errorf("audio = %q", got)
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

// TestSpeakRespells checks that a provider that ignores Pronunciations
// gets the respellings in its text.
func TestSpeakRespells(t *testing.T) {
	f := registerFake(t)
	delete(appliesPronunciations, "fake")
	o := &speakOptions{provider: "fake", voice: "v", output: filepath.Join(t.TempDir(), "a.mp3"), phonemes: "auto",
		ipa: []string{"Miura=ˈmjura", "Countach=ˈkuntaʃ"}, say: []string{"Miura=mee oora"}}
	if err := runSpeak(context.Background(), o, []string{"The Miura and the Countach."}, nil); err != nil {
		t.Fatal(err)
	}
	if f.text != "The mee oora and the Countach." || f.config.Pronunciations != nil {
		t.Errorf("request = %q, %+v", f.text, f.config.Pronunciations)
	}
}

func TestAudioFile(t *testing.T) {
	raw := []byte{1, 2, 3, 4}
	cases := []struct {
		name, format string
		res          tts.SynthesisResult
		want         string // "raw", "wrapped" or "error"
	}{
		{"mp3", "mp3", tts.SynthesisResult{Audio: raw}, "raw"},
		{"pcm", "pcm", tts.SynthesisResult{Audio: raw}, "raw"},
		{"wav from the provider", "wav", tts.SynthesisResult{Audio: []byte("RIFF\x00\x00\x00\x00WAVEfmt ")}, "raw"},
		{"raw PCM for a wav", "wav", tts.SynthesisResult{Audio: raw, SampleRate: 44100}, "wrapped"},
		{"raw PCM without a rate", "wav", tts.SynthesisResult{Audio: raw}, "error"},
		{"ogg opus", "opus", tts.SynthesisResult{Audio: []byte("OggS....")}, "raw"},
		{"mp3 for an opus", "opus", tts.SynthesisResult{Audio: []byte("ID3....")}, "error"},
	}
	for _, tc := range cases {
		got, err := audioFile(tc.format, &tc.res)
		switch {
		case tc.want == "error":
			if err == nil {
				t.Errorf("%s: no error", tc.name)
			}
		case err != nil:
			t.Errorf("%s: %v", tc.name, err)
		case tc.want == "raw" && !bytes.Equal(got, tc.res.Audio):
			t.Errorf("%s: audio changed", tc.name)
		case tc.want == "wrapped":
			if len(got) != 44+len(raw) || string(got[:4]) != "RIFF" || string(got[8:16]) != "WAVEfmt " ||
				binary.LittleEndian.Uint32(got[24:]) != 44100 || binary.LittleEndian.Uint32(got[40:]) != 4 ||
				!bytes.Equal(got[44:], raw) {
				t.Errorf("%s: header % x", tc.name, got[:44])
			}
		}
	}
}

// TestTTSProviderKey checks that the API key is optional: the provider's
// factory decides, and its error names the variable to set.
func TestTTSProviderKey(t *testing.T) {
	omnivoice.RegisterTTSProvider("keyless", func(omnivoice.ProviderConfig) (tts.Provider, error) { return &fakeTTS{}, nil }, omnivoice.PriorityThick)
	omnivoice.RegisterTTSProvider("keyed", func(c omnivoice.ProviderConfig) (tts.Provider, error) {
		if c.APIKey == "" {
			return nil, errors.New("keyed: API key is required")
		}
		return &fakeTTS{}, nil
	}, omnivoice.PriorityThick)
	t.Setenv("KEYLESS_API_KEY", "")
	t.Setenv("KEYED_API_KEY", "")
	if _, err := ttsProvider("keyless"); err != nil {
		t.Errorf("keyless: %v", err)
	}
	if _, err := ttsProvider("keyed"); err == nil || !strings.Contains(err.Error(), "KEYED_API_KEY") {
		t.Errorf("keyed without a key: %v", err)
	}
	t.Setenv("KEYED_API_KEY", "k")
	if _, err := ttsProvider("keyed"); err != nil {
		t.Errorf("keyed with a key: %v", err)
	}
}
