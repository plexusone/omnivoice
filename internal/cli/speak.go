package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/plexusone/omnivoice"
	"github.com/plexusone/omnivoice-core/terminology"
	"github.com/spf13/cobra"
)

type speakOptions struct {
	provider string
	voice    string
	model    string
	output   string
	locale   string
	termsDir string
	ipa      []string
	say      []string
	phonemes string
}

func newSpeakCmd() *cobra.Command {
	o := &speakOptions{}
	cmd := &cobra.Command{
		Use:   "speak [text]",
		Short: "Synthesize speech from text",
		Long: `Synthesize speech from text with the chosen provider and write the audio
to a file. Read the text from standard input with "-" or no argument.

Pronunciations for names and terms come from a terminology-spec terms
directory (--terms), and from --ipa and --say for quick trials. Providers
use a term's IPA where the voice or model accepts phonemes, and its
respelling otherwise.`,
		Example: `  omnivoice speak --voice <id> -o miura.mp3 "The Miura in sixty-six."
  omnivoice speak --voice <id> --model eleven_flash_v2 --ipa "Miura=ˈmjura" -o miura.mp3 "The Miura."
  omnivoice speak --voice <id> --terms terms -o lesson.mp3 - < lesson.txt`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSpeak(cmd.Context(), o, args, cmd.InOrStdin())
		},
	}
	f := cmd.Flags()
	f.StringVarP(&o.provider, "provider", "p", "elevenlabs", "TTS provider")
	f.StringVar(&o.voice, "voice", "", "voice ID (list them with: omnivoice voices)")
	f.StringVarP(&o.model, "model", "m", "", "provider-specific model")
	f.StringVarP(&o.output, "output", "o", "", "audio file to write (required)")
	f.StringVarP(&o.locale, "locale", "l", "en-US", "BCP-47 locale of the text, for choosing pronunciations")
	f.StringVar(&o.termsDir, "terms", "", "terminology-spec terms directory")
	f.StringArrayVar(&o.ipa, "ipa", nil, "TERM=IPA, such as Miura=ˈmjura (repeatable)")
	f.StringArrayVar(&o.say, "say", nil, "TERM=RESPELLING, such as Miura=mee oora (repeatable)")
	f.StringVar(&o.phonemes, "phonemes", "auto", "ElevenLabs: send IPA as phoneme tags: auto (by model), on or off")
	if err := cmd.MarkFlagRequired("voice"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("output"); err != nil {
		panic(err)
	}
	return cmd
}

func runSpeak(ctx context.Context, o *speakOptions, args []string, stdin io.Reader) error {
	text, err := speakText(args, stdin)
	if err != nil {
		return err
	}
	pronunciations, err := o.pronunciations()
	if err != nil {
		return err
	}
	config := omnivoice.SynthesisConfig{
		VoiceID:        o.voice,
		Model:          o.model,
		OutputFormat:   outputFormat(o.output),
		Pronunciations: pronunciations,
	}
	switch o.phonemes {
	case "auto":
	case "on", "off":
		config.Extensions = map[string]any{"elevenlabs.phonemes": o.phonemes == "on"}
	default:
		return fmt.Errorf("--phonemes must be auto, on or off, not %q", o.phonemes)
	}

	p, err := ttsProvider(o.provider)
	if err != nil {
		return err
	}
	logVerbose("Synthesizing %d characters with %s, %d pronunciations", len(text), o.provider, len(pronunciations))
	res, err := p.Synthesize(ctx, text, config)
	if err != nil {
		return fmt.Errorf("synthesis failed: %w", err)
	}
	if err := os.WriteFile(o.output, res.Audio, 0o600); err != nil {
		return err
	}
	logInfo("Wrote %s (%d bytes)", o.output, len(res.Audio))
	return nil
}

func speakText(args []string, stdin io.Reader) (string, error) {
	if len(args) == 1 && args[0] != "-" {
		return args[0], nil
	}
	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", err
	}
	text := strings.TrimSpace(string(data))
	if text == "" {
		return "", errors.New("no text to speak")
	}
	return text, nil
}

// pronunciations merges the terms directory with --ipa and --say, which
// win for the terms they name.
func (o *speakOptions) pronunciations() ([]omnivoice.Pronunciation, error) {
	var ps []omnivoice.Pronunciation
	if o.termsDir != "" {
		terms, err := terminology.LoadDir(o.termsDir)
		if err != nil {
			return nil, fmt.Errorf("load terms: %w", err)
		}
		ps = terminology.PronunciationsFor(terms, o.locale)
	}
	byTerm := map[string]int{}
	for i, p := range ps {
		byTerm[strings.ToLower(p.Term)] = i
	}
	set := func(flag, value string, apply func(*omnivoice.Pronunciation, string)) error {
		term, v, ok := strings.Cut(value, "=")
		term, v = strings.TrimSpace(term), strings.TrimSpace(v)
		if !ok || term == "" || v == "" {
			return fmt.Errorf("--%s %q: want TERM=VALUE", flag, value)
		}
		i, seen := byTerm[strings.ToLower(term)]
		if !seen {
			ps = append(ps, omnivoice.Pronunciation{Term: term})
			i = len(ps) - 1
			byTerm[strings.ToLower(term)] = i
		}
		apply(&ps[i], v)
		return nil
	}
	for _, v := range o.ipa {
		if err := set("ipa", v, func(p *omnivoice.Pronunciation, s string) { p.IPA = strings.Trim(s, "/[]") }); err != nil {
			return nil, err
		}
	}
	for _, v := range o.say {
		if err := set("say", v, func(p *omnivoice.Pronunciation, s string) { p.Respelling = s }); err != nil {
			return nil, err
		}
	}
	return ps, nil
}

// outputFormat names the audio format from the output file's extension.
func outputFormat(path string) string {
	i := strings.LastIndex(path, ".")
	if i < 0 {
		return "mp3"
	}
	switch ext := strings.ToLower(path[i+1:]); ext {
	case "wav", "pcm", "opus", "mp3":
		return ext
	default:
		return "mp3"
	}
}

func ttsProvider(name string) (omnivoice.TTSProvider, error) {
	key := getAPIKeyForProvider(name)
	if key == "" {
		return nil, fmt.Errorf("set %s for %s", getEnvVarName(name), name)
	}
	return omnivoice.GetTTSProvider(name, omnivoice.WithAPIKey(key))
}

type voicesOptions struct {
	provider string
	match    string
	json     bool
}

func newVoicesCmd() *cobra.Command {
	o := &voicesOptions{}
	cmd := &cobra.Command{
		Use:   "voices",
		Short: "List a TTS provider's voices",
		Example: `  omnivoice voices --match italian
  omnivoice voices -p openai --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVoices(cmd.Context(), o, cmd.OutOrStdout())
		},
	}
	cmd.Flags().StringVarP(&o.provider, "provider", "p", "elevenlabs", "TTS provider")
	cmd.Flags().StringVar(&o.match, "match", "", "only voices whose name, language or details contain this text")
	cmd.Flags().BoolVar(&o.json, "json", false, "print JSON")
	return cmd
}

func runVoices(ctx context.Context, o *voicesOptions, w io.Writer) error {
	p, err := ttsProvider(o.provider)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	vs, err := p.ListVoices(ctx)
	if err != nil {
		return err
	}
	sort.Slice(vs, func(i, j int) bool { return vs[i].Name < vs[j].Name })
	match := strings.ToLower(o.match)
	var shown []omnivoice.Voice
	for _, v := range vs {
		meta, err := json.Marshal(v.Metadata)
		if err != nil {
			return err
		}
		if match == "" || strings.Contains(strings.ToLower(v.Name+" "+v.Language+" "+string(meta)), match) {
			shown = append(shown, v)
		}
	}
	if o.json {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(shown)
	}
	for _, v := range shown {
		accent := ""
		if a, ok := v.Metadata["accent"].(string); ok {
			accent = a
		}
		if _, err := fmt.Fprintf(w, "%-24s %-40s %-6s %-8s %s\n", v.ID, v.Name, v.Language, v.Gender, accent); err != nil {
			return err
		}
	}
	return nil
}
