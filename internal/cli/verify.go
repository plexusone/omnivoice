package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/grokify/mogo/time/duration"
	"github.com/plexusone/omnivoice"
	"github.com/plexusone/omnivoice-core/stt"
	"github.com/plexusone/omnivoice-core/stt/verify"
	"github.com/plexusone/omnivoice-core/terminology"
	"github.com/spf13/cobra"
)

func newVerifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Check quotations against a recording",
		Long: `Check wording attributed to a recording by comparing what several
speech-to-text engines hear in a clip of it.

Each quotation is located in an index transcript, cut out with a little
padding, and transcribed by every engine. The words are normalized and
compared, and each quotation gets a decision:

  agreed        every engine heard the quotation
  corrected     the engines agree with each other but not with the quotation
  disputed      the engines disagree on a number, negation, name, term or
                other substantive word
  notFound      the quotation could not be located in the recording
  insufficient  too few engines produced a transcript

Names, titles and dates should be confirmed from a written source, whatever
the engines hear.`,
	}
	cmd.AddCommand(newVerifyQuotesCmd(), newVerifyAbsenceCmd())
	return cmd
}

type verifyQuotesOptions struct {
	audio       string
	index       string
	quotesFile  string
	quotes      []string
	providers   string
	models      []string
	keywords    []string
	termsDir    string
	plainEngine string
	minEngines  int
	pad         time.Duration
	language    string
	output      string
}

func newVerifyQuotesCmd() *cobra.Command {
	o := &verifyQuotesOptions{}
	cmd := &cobra.Command{
		Use:   "quotes",
		Short: "Check quotations against a recording with several engines",
		Long: `Check quotations against a recording with several engines.

Quotations come from --quotes, a JSON array of {"id", "text", "start_ms",
"end_ms"} where the times are optional, or from --quote, repeated. A
quotation without times is located in the --index transcript, an OmniVoice
JSON transcript of the whole recording with word timestamps.

API keys come from the environment: DEEPGRAM_API_KEY, OPENAI_API_KEY and
ELEVENLABS_API_KEY. ffmpeg must be on PATH.

Key words (--keyword, and the terms in --terms) help engines with names and
terms. They are never taken from the quotation, and one engine always runs
without them (--plain-engine, by default the last in --providers).

Examples:
  omnivoice verify quotes --audio ep753.mp3 --index transcript.json \
    --quotes quotes.json -p deepgram,elevenlabs,openai -o report.json

  omnivoice verify quotes --audio ep753.mp3 --index transcript.json \
    --quote "we can have the model run all the tests" \
    -p deepgram,elevenlabs --min-engines 2 --keyword Mantle --keyword Bedrock`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVerifyQuotes(cmd.Context(), o, cmd.OutOrStdout())
		},
	}
	f := cmd.Flags()
	f.StringVar(&o.audio, "audio", "", "recording to check against (required)")
	f.StringVar(&o.index, "index", "", "OmniVoice JSON transcript of the recording, with word timestamps")
	f.StringVar(&o.quotesFile, "quotes", "", "JSON file of quotations")
	f.StringArrayVar(&o.quotes, "quote", nil, "a quotation to check; repeat for more")
	f.StringVarP(&o.providers, "providers", "p", "deepgram,elevenlabs,openai", "comma-separated engines")
	f.StringArrayVar(&o.models, "model", nil, "engine model as provider=model; repeat for more")
	f.StringArrayVar(&o.keywords, "keyword", nil, "a name or term to help engines; repeat for more")
	f.StringVar(&o.termsDir, "terms", "", "terminology-spec terms directory: canonical forms become key words, aliases are accepted")
	f.StringVar(&o.plainEngine, "plain-engine", "", "engine that runs without key words (default: the last in --providers)")
	f.IntVar(&o.minEngines, "min-engines", 3, "fewest engines that must produce a transcript to decide")
	f.DurationVar(&o.pad, "pad", 1500*time.Millisecond, "padding before and after each quotation's clip")
	f.StringVarP(&o.language, "language", "l", "en-US", "language code")
	f.StringVarP(&o.output, "output", "o", "", "write the JSON report here (default: stdout)")
	if err := cmd.MarkFlagRequired("audio"); err != nil {
		panic(err) // the flag is defined just above
	}
	return cmd
}

func runVerifyQuotes(ctx context.Context, o *verifyQuotesOptions, stdout io.Writer) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, err := os.Stat(o.audio); err != nil {
		return fmt.Errorf("audio: %w", err)
	}
	quotes, err := o.loadQuotes()
	if err != nil {
		return err
	}
	if len(quotes) == 0 {
		return fmt.Errorf("no quotations: use --quotes or --quote")
	}
	var index []verify.TimedWord
	indexName := ""
	if o.index != "" {
		t, err := stt.LoadTranscript(o.index)
		if err != nil {
			return fmt.Errorf("index: %w", err)
		}
		index = verify.WordsFromTranscript(t)
		indexName = strings.Trim(t.Metadata.Provider+"/"+t.Metadata.Model, "/")
		if len(index) == 0 {
			return fmt.Errorf("index %s has no word timestamps", o.index)
		}
	}

	keywords := append([]string(nil), o.keywords...)
	var termNames []string
	aliases := map[string][]string{}
	if o.termsDir != "" {
		terms, err := terminology.LoadDir(o.termsDir)
		if err != nil {
			return fmt.Errorf("terms: %w", err)
		}
		for _, t := range terms {
			if t.CanonicalForm == "" {
				continue
			}
			keywords = append(keywords, t.CanonicalForm)
			termNames = append(termNames, t.CanonicalForm)
			if len(t.Aliases) > 0 {
				aliases[t.CanonicalForm] = append(aliases[t.CanonicalForm], t.Aliases...)
			}
		}
	}
	termNames = append(termNames, o.keywords...)

	engines, err := o.engines(len(keywords) > 0)
	if err != nil {
		return err
	}
	// Fail before any engine runs: too few engines can never decide.
	if o.minEngines < 1 || o.minEngines > len(engines) {
		return fmt.Errorf("--min-engines %d needs between 1 and %d with %d provider(s) (%s); pass --min-engines %d or add providers",
			o.minEngines, len(engines), len(engines), o.providers, len(engines))
	}
	v := verify.Verifier{
		Engines: engines,
		Options: verify.Options{
			Pad:      duration.FromDuration(o.pad),
			Keywords: keywords,
			Compare: verify.CompareOptions{
				MinEngines: o.minEngines,
				Normalize:  verify.NormalizeOptions{Terms: termNames},
				Aliases:    aliases,
			},
		},
	}
	logInfo("Checking %d quotation(s) with %s...", len(quotes), o.providers)
	rep, err := v.VerifyQuotes(ctx, o.audio, index, indexName, quotes)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if o.output != "" {
		if err := os.WriteFile(o.output, data, 0o600); err != nil {
			return fmt.Errorf("write report: %w", err)
		}
		logInfo("Report written to %s", o.output)
	} else if _, err := stdout.Write(data); err != nil {
		return err
	}
	printSummary(rep)
	return nil
}

func (o *verifyQuotesOptions) loadQuotes() ([]verify.Quote, error) {
	var quotes []verify.Quote
	if o.quotesFile != "" {
		b, err := os.ReadFile(o.quotesFile)
		if err != nil {
			return nil, fmt.Errorf("quotes: %w", err)
		}
		if err := json.Unmarshal(b, &quotes); err != nil {
			return nil, fmt.Errorf("quotes %s: %w", o.quotesFile, err)
		}
	}
	for i, q := range o.quotes {
		quotes = append(quotes, verify.Quote{ID: fmt.Sprintf("quote-%d", i+1), Text: q})
	}
	seen := map[string]bool{}
	for i, q := range quotes {
		if q.ID == "" {
			return nil, fmt.Errorf("quotation %d has no id", i+1)
		}
		if seen[q.ID] {
			return nil, fmt.Errorf("quotation id %q is used twice", q.ID)
		}
		seen[q.ID] = true
		if q.Start == 0 && q.End == 0 && o.index == "" {
			return nil, fmt.Errorf("quotation %q has no times; give --index to locate it", q.ID)
		}
	}
	return quotes, nil
}

// engines builds one engine per provider. When key words are in use, the
// plain engine runs without them.
func (o *verifyQuotesOptions) engines(withKeywords bool) ([]verify.Engine, error) {
	models := map[string]string{}
	for _, m := range o.models {
		p, model, ok := strings.Cut(m, "=")
		if !ok || p == "" || model == "" {
			return nil, fmt.Errorf("--model %q: want provider=model", m)
		}
		models[strings.ToLower(p)] = model
	}
	var names []string
	seen := map[string]bool{}
	for _, n := range strings.Split(o.providers, ",") {
		n = strings.ToLower(strings.TrimSpace(n))
		if n == "" {
			continue
		}
		if seen[n] {
			// The same engine twice would count as two engines agreeing.
			return nil, fmt.Errorf("--providers lists %s more than once", n)
		}
		seen[n] = true
		names = append(names, n)
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("no providers")
	}
	plain := strings.ToLower(o.plainEngine)
	if plain == "" {
		plain = names[len(names)-1]
	}
	var engines []verify.Engine
	foundPlain := false
	for _, n := range names {
		key := getAPIKeyForProvider(n)
		if key == "" {
			return nil, fmt.Errorf("API key not found for %s (set %s)", n, getEnvVarName(n))
		}
		p, err := omnivoice.GetSTTProvider(n, omnivoice.WithAPIKey(key))
		if err != nil {
			return nil, fmt.Errorf("provider %s: %w", n, err)
		}
		e := verify.Engine{
			Name:     n,
			Provider: p,
			Config: stt.TranscriptionConfig{
				Language:          o.language,
				Model:             models[n],
				EnablePunctuation: true,
			},
			NoKeywords: n == plain,
		}
		if e.NoKeywords {
			foundPlain = true
		}
		engines = append(engines, e)
	}
	if withKeywords && !foundPlain {
		return nil, fmt.Errorf("--plain-engine %q is not in --providers", o.plainEngine)
	}
	return engines, nil
}

func printSummary(rep *verify.Report) {
	if quiet {
		return
	}
	fmt.Fprintln(os.Stderr)
	for _, r := range rep.Results {
		c := r.Comparison
		fmt.Fprintf(os.Stderr, "%-12s %-12s %s–%s  %s\n", r.Quote.ID, c.Decision,
			clock(r.ClipStart.Duration()), clock(r.ClipEnd.Duration()), short(r.Quote.Text, 60))
		if c.Corrected != "" {
			fmt.Fprintf(os.Stderr, "%26sheard: %s\n", "", short(c.Corrected, 80))
		}
		for _, d := range c.Differences {
			if !d.Substantive {
				continue
			}
			if d.Insertion {
				fmt.Fprintf(os.Stderr, "%26sextra words (%s): %s\n", "", d.Kind, heardList(d.Heard))
				continue
			}
			fmt.Fprintf(os.Stderr, "%26s%s %q: %s\n", "", d.Kind, d.Expected, heardList(d.Heard))
		}
		if r.Err != "" {
			fmt.Fprintf(os.Stderr, "%26serror: %s\n", "", r.Err)
		}
	}
	keys := make([]string, 0, len(rep.Counts))
	for k := range rep.Counts {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s %d", k, rep.Counts[verify.Decision(k)]))
	}
	fmt.Fprintf(os.Stderr, "\n%s\n", strings.Join(parts, ", "))
}

func heardList(h map[string]string) string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		v := h[k]
		if v == "" {
			v = "—"
		}
		parts = append(parts, k+" "+v)
	}
	return strings.Join(parts, ", ")
}

func clock(d time.Duration) string {
	s := int(d.Seconds())
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

func short(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}

type verifyAbsenceOptions struct {
	transcripts []string
	phrases     []string
	minScore    float64
	output      string
}

func newVerifyAbsenceCmd() *cobra.Command {
	o := &verifyAbsenceOptions{}
	cmd := &cobra.Command{
		Use:   "absence",
		Short: "Search independent transcripts for phrases a recording may not contain",
		Long: `Search several independent transcripts of one recording for phrases.

Use it before saying a recording doesn't mention something. The result
never says a phrase is absent from the recording, only whether it was found
in each transcript. Give transcripts from different engines: one engine can
miss a passage.

Example:
  omnivoice verify absence --transcript deepgram.json --transcript whisper.json \
    --phrase "locally" --phrase "on my laptop" --phrase "run the whole stack"`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVerifyAbsence(o, cmd.OutOrStdout())
		},
	}
	f := cmd.Flags()
	f.StringArrayVar(&o.transcripts, "transcript", nil, "OmniVoice JSON transcript with word timestamps; repeat for each engine")
	f.StringArrayVar(&o.phrases, "phrase", nil, "phrase to search for; repeat for more")
	f.Float64Var(&o.minScore, "min-score", 0.8, "lowest match score that counts as found (0–1)")
	f.StringVarP(&o.output, "output", "o", "", "write the JSON result here (default: stdout)")
	return cmd
}

func runVerifyAbsence(o *verifyAbsenceOptions, stdout io.Writer) error {
	if len(o.transcripts) == 0 || len(o.phrases) == 0 {
		return fmt.Errorf("give at least one --transcript and one --phrase")
	}
	if o.minScore < 0 || o.minScore > 1 {
		return fmt.Errorf("--min-score %v is outside 0–1 (for example, 0.8)", o.minScore)
	}
	var ts []verify.IndexedTranscript
	for _, p := range o.transcripts {
		t, err := stt.LoadTranscript(p)
		if err != nil {
			return fmt.Errorf("transcript %s: %w", p, err)
		}
		words := verify.WordsFromTranscript(t)
		if len(words) == 0 {
			return fmt.Errorf("transcript %s has no word timestamps", p)
		}
		engine := t.Metadata.Provider
		if engine == "" {
			engine = p
		}
		ts = append(ts, verify.IndexedTranscript{Engine: engine, Model: t.Metadata.Model, Words: words})
	}
	if len(ts) < 2 {
		logInfo("Only one transcript: a phrase it lacks may still be in the recording.")
	}
	res := verify.Search(ts, o.phrases, verify.AbsenceOptions{MinScore: o.minScore})
	data, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if o.output != "" {
		if err := os.WriteFile(o.output, data, 0o600); err != nil {
			return fmt.Errorf("write result: %w", err)
		}
	} else if _, err := stdout.Write(data); err != nil {
		return err
	}
	if !quiet {
		for _, p := range res.Phrases {
			fmt.Fprintf(os.Stderr, "%-40s found in %d of %d\n", short(p.Phrase, 40), p.FoundIn, res.Transcripts)
		}
		fmt.Fprintln(os.Stderr, res.Statement)
	}
	return nil
}
