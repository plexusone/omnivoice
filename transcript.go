package omnivoice

import (
	"github.com/plexusone/omnivoice-core/stt"
)

// Re-export Transcript types from omnivoice-core/stt for backwards compatibility.
// The canonical Transcript format is now defined in omnivoice-core.

// TranscriptFormatVersion is the current version of the OmniVoice transcript format.
const TranscriptFormatVersion = stt.TranscriptFormatVersion

// TranscriptSchemaURL is the JSON Schema URL for the transcript format.
const TranscriptSchemaURL = stt.TranscriptSchemaURL

// Type aliases for backwards compatibility.
type (
	Transcript         = stt.Transcript
	TranscriptSegment  = stt.TranscriptSegment
	TranscriptWord     = stt.TranscriptWord
	TranscriptMetadata = stt.TranscriptMetadata
	TranscriptOptions  = stt.TranscriptOptions
)

// NewTranscript creates a Transcript from a TranscriptionResult.
// This is a convenience wrapper around stt.NewTranscript.
func NewTranscript(result *TranscriptionResult, provider, model, audioFile string, config *TranscriptionConfig) *Transcript {
	// Convert omnivoice.TranscriptionConfig to stt.TranscriptionConfig
	var sttConfig *stt.TranscriptionConfig
	if config != nil {
		sttConfig = &stt.TranscriptionConfig{
			Language:                 config.Language,
			EnablePunctuation:        config.EnablePunctuation,
			EnableWordTimestamps:     config.EnableWordTimestamps,
			EnableSpeakerDiarization: config.EnableSpeakerDiarization,
		}
	}

	// Convert omnivoice.TranscriptionResult to stt.TranscriptionResult
	sttResult := &stt.TranscriptionResult{
		Text:               result.Text,
		Language:           result.Language,
		LanguageConfidence: result.LanguageConfidence,
		Duration:           result.Duration,
	}

	// Convert segments
	sttResult.Segments = make([]stt.Segment, len(result.Segments))
	for i, seg := range result.Segments {
		sttResult.Segments[i] = stt.Segment{
			Text:       seg.Text,
			StartTime:  seg.StartTime,
			EndTime:    seg.EndTime,
			Confidence: seg.Confidence,
			Speaker:    seg.Speaker,
			Language:   seg.Language,
		}
		sttResult.Segments[i].Words = make([]stt.Word, len(seg.Words))
		for j, word := range seg.Words {
			sttResult.Segments[i].Words[j] = stt.Word{
				Text:       word.Text,
				StartTime:  word.StartTime,
				EndTime:    word.EndTime,
				Confidence: word.Confidence,
				Speaker:    word.Speaker,
			}
		}
	}

	return stt.NewTranscript(sttResult, provider, model, audioFile, sttConfig)
}

// LoadTranscript reads a transcript from a JSON file.
func LoadTranscript(filePath string) (*Transcript, error) {
	return stt.LoadTranscript(filePath)
}
