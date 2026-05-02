package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/plexusone/omnivoice"
	_ "github.com/plexusone/omnivoice/providers/all"
	"github.com/spf13/cobra"
)

var (
	provider   string
	outputFile string
	format     string
	language   string
	diarize    bool
	timestamps bool
	model      string
)

var transcribeCmd = &cobra.Command{
	Use:   "transcribe <file>",
	Short: "Transcribe audio file to text",
	Long: `Transcribe an audio file to text using the specified provider.

Supported audio formats: mp3, wav, flac, ogg, m4a, webm

Examples:
  # Basic transcription (stdout)
  omnivoice transcribe podcast.mp3

  # With specific provider and output file
  omnivoice transcribe -p deepgram -o transcript.txt podcast.mp3

  # JSON output with full metadata (OmniVoice Transcript format)
  omnivoice transcribe -p deepgram --diarize --timestamps -f json -o transcript.json podcast.mp3

  # Generate SRT subtitles
  omnivoice transcribe -p deepgram -f srt -o subtitles.srt podcast.mp3

  # Generate WebVTT subtitles (for web video)
  omnivoice transcribe -p deepgram -f vtt -o subtitles.vtt podcast.mp3

  # With speaker diarization
  omnivoice transcribe -p deepgram --diarize -f json -o transcript.json meeting.mp3`,
	Args: cobra.ExactArgs(1),
	RunE: runTranscribe,
}

func init() {
	transcribeCmd.Flags().StringVarP(&provider, "provider", "p", "deepgram", "STT provider (deepgram, openai, elevenlabs)")
	transcribeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path (default: stdout)")
	transcribeCmd.Flags().StringVarP(&format, "format", "f", "text", "output format: text, json, srt, vtt")
	transcribeCmd.Flags().StringVarP(&language, "language", "l", "en-US", "language code (e.g., en-US)")
	transcribeCmd.Flags().BoolVar(&diarize, "diarize", false, "enable speaker diarization")
	transcribeCmd.Flags().BoolVar(&timestamps, "timestamps", false, "enable word timestamps")
	transcribeCmd.Flags().StringVarP(&model, "model", "m", "", "provider-specific model")
}

func runTranscribe(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Get API key from environment
	apiKey := getAPIKeyForProvider(provider)
	if apiKey == "" {
		return fmt.Errorf("API key not found for provider %s (set %s environment variable)", provider, getEnvVarName(provider))
	}

	// Create provider
	logVerbose("Using provider: %s", provider)
	sttProvider, err := omnivoice.GetSTTProvider(provider, omnivoice.WithAPIKey(apiKey))
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	// Configure transcription
	config := omnivoice.TranscriptionConfig{
		Language:                 language,
		Model:                    model,
		EnablePunctuation:        true,
		EnableWordTimestamps:     timestamps || format == "srt" || format == "vtt",
		EnableSpeakerDiarization: diarize,
	}

	// Transcribe
	logInfo("Transcribing %s...", filePath)
	ctx := context.Background()
	result, err := sttProvider.TranscribeFile(ctx, filePath, config)
	if err != nil {
		return fmt.Errorf("transcription failed: %w", err)
	}

	logVerbose("Duration: %s", result.Duration)
	logVerbose("Language: %s (confidence: %.2f)", result.Language, result.LanguageConfidence)

	// Format output
	output, err := formatOutput(result, format, provider, model, filePath, &config)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	// Write output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(output), 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		logInfo("Output written to %s", outputFile)
	} else {
		fmt.Print(output)
	}

	return nil
}

func getAPIKeyForProvider(provider string) string {
	switch strings.ToLower(provider) {
	case "deepgram":
		return os.Getenv("DEEPGRAM_API_KEY")
	case "openai":
		return os.Getenv("OPENAI_API_KEY")
	case "elevenlabs":
		return os.Getenv("ELEVENLABS_API_KEY")
	default:
		// Try provider name as uppercase with _API_KEY suffix
		return os.Getenv(strings.ToUpper(provider) + "_API_KEY")
	}
}

func getEnvVarName(provider string) string {
	switch strings.ToLower(provider) {
	case "deepgram":
		return "DEEPGRAM_API_KEY"
	case "openai":
		return "OPENAI_API_KEY"
	case "elevenlabs":
		return "ELEVENLABS_API_KEY"
	default:
		return strings.ToUpper(provider) + "_API_KEY"
	}
}

func formatOutput(result *omnivoice.TranscriptionResult, format, providerName, modelName, audioFile string, config *omnivoice.TranscriptionConfig) (string, error) {
	switch strings.ToLower(format) {
	case "text":
		return result.Text + "\n", nil

	case "json":
		transcript := omnivoice.NewTranscript(result, providerName, modelName, audioFile, config)
		data, err := transcript.ToJSON()
		if err != nil {
			return "", err
		}
		return string(data) + "\n", nil

	case "srt":
		opts := omnivoice.DefaultSubtitleOptions()
		opts.IncludeSpeakerLabels = diarize
		return omnivoice.GenerateSRT(result, opts), nil

	case "vtt":
		opts := omnivoice.DefaultSubtitleOptions()
		opts.IncludeSpeakerLabels = diarize
		return omnivoice.GenerateVTT(result, opts), nil

	default:
		return "", fmt.Errorf("unknown format: %s (supported: text, json, srt, vtt)", format)
	}
}
