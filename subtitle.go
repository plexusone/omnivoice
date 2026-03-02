package omnivoice

import (
	"github.com/plexusone/omnivoice-core/subtitle"
)

// Re-export subtitle types from omnivoice-core
type (
	// SubtitleOptions configures subtitle generation.
	SubtitleOptions = subtitle.Options

	// SubtitleFormat represents the output format for subtitles.
	SubtitleFormat = subtitle.Format
)

// Subtitle format constants.
const (
	SubtitleFormatSRT = subtitle.FormatSRT
	SubtitleFormatVTT = subtitle.FormatVTT
)

// Re-export subtitle functions.
var (
	// DefaultSubtitleOptions returns sensible defaults for subtitle generation.
	DefaultSubtitleOptions = subtitle.DefaultOptions

	// GenerateSRT generates SRT subtitles from a transcription result.
	GenerateSRT = subtitle.GenerateSRT

	// GenerateVTT generates WebVTT subtitles from a transcription result.
	GenerateVTT = subtitle.GenerateVTT

	// SaveSRT generates and saves SRT to a file.
	SaveSRT = subtitle.SaveSRT

	// SaveVTT generates and saves WebVTT to a file.
	SaveVTT = subtitle.SaveVTT
)
