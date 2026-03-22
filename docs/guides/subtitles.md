# Subtitles

Generate SRT and VTT subtitle files from transcription results.

## Quick Start

```go
import "github.com/plexusone/omnivoice"

// Transcribe with word timestamps
result, _ := sttProvider.TranscribeFile(ctx, "video.mp4", omnivoice.TranscriptionConfig{
    EnableWordTimestamps: true,
})

// Generate SRT subtitles
srt := omnivoice.GenerateSRT(result)
os.WriteFile("video.srt", []byte(srt), 0644)

// Generate VTT subtitles
vtt := omnivoice.GenerateVTT(result)
os.WriteFile("video.vtt", []byte(vtt), 0644)
```

## Subtitle Formats

### SRT (SubRip)

```
1
00:00:00,000 --> 00:00:03,500
Hello and welcome to our video.

2
00:00:03,500 --> 00:00:07,200
Today we'll be discussing
speech recognition.
```

### VTT (WebVTT)

```
WEBVTT

00:00:00.000 --> 00:00:03.500
Hello and welcome to our video.

00:00:03.500 --> 00:00:07.200
Today we'll be discussing
speech recognition.
```

## Configuration

### SubtitleConfig

```go
config := omnivoice.SubtitleConfig{
    MaxCharsPerLine: 42,    // Characters per line
    MaxLinesPerCue:  2,     // Lines per subtitle cue
    MinCueDuration:  1.0,   // Minimum cue duration (seconds)
    MaxCueDuration:  7.0,   // Maximum cue duration (seconds)
}

srt := omnivoice.GenerateSRTWithConfig(result, config)
```

### Default Values

| Setting | Default | Description |
|---------|---------|-------------|
| MaxCharsPerLine | 42 | Standard for readability |
| MaxLinesPerCue | 2 | Standard subtitle format |
| MinCueDuration | 1.0s | Minimum display time |
| MaxCueDuration | 7.0s | Maximum display time |

## Transcription Requirements

Word timestamps are required for accurate subtitle timing:

```go
result, _ := provider.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
    Language:             "en",
    EnableWordTimestamps: true, // Required for subtitles
})

// Check if word timestamps are available
if len(result.Words) == 0 {
    log.Println("Warning: No word timestamps, using segment-based timing")
}
```

## Complete Example

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"
)

func main() {
    ctx := context.Background()

    // Get STT provider
    provider, err := omnivoice.GetSTTProvider("deepgram",
        omnivoice.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))
    if err != nil {
        log.Fatal(err)
    }

    // Transcribe with word timestamps
    result, err := provider.TranscribeFile(ctx, "video.mp4", omnivoice.TranscriptionConfig{
        Language:             "en",
        EnableWordTimestamps: true,
        Extensions: map[string]any{
            "punctuate":    true,
            "smart_format": true,
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Generate subtitles with custom config
    config := omnivoice.SubtitleConfig{
        MaxCharsPerLine: 40,
        MaxLinesPerCue:  2,
    }

    // Save SRT
    srt := omnivoice.GenerateSRTWithConfig(result, config)
    if err := os.WriteFile("video.srt", []byte(srt), 0644); err != nil {
        log.Fatal(err)
    }

    // Save VTT
    vtt := omnivoice.GenerateVTTWithConfig(result, config)
    if err := os.WriteFile("video.vtt", []byte(vtt), 0644); err != nil {
        log.Fatal(err)
    }

    log.Println("Subtitles generated: video.srt, video.vtt")
}
```

## Handling Different Languages

### Chinese/Japanese (No Spaces)

For languages without word boundaries, the subtitle generator handles character-by-character grouping:

```go
result, _ := provider.TranscribeFile(ctx, "chinese_audio.mp3", omnivoice.TranscriptionConfig{
    Language:             "zh",
    EnableWordTimestamps: true,
})

// Generator handles character grouping automatically
srt := omnivoice.GenerateSRT(result)
```

### Mixed Language Content

```go
result, _ := provider.TranscribeFile(ctx, "mixed_audio.mp3", omnivoice.TranscriptionConfig{
    Language:             "en", // Primary language
    EnableWordTimestamps: true,
    Extensions: map[string]any{
        "detect_language": true,
    },
})
```

## Speaker Diarization

Include speaker labels in subtitles:

```go
result, _ := provider.TranscribeFile(ctx, "interview.mp3", omnivoice.TranscriptionConfig{
    EnableWordTimestamps:     true,
    EnableSpeakerDiarization: true,
})

// Generate subtitles with speaker labels
for _, segment := range result.Segments {
    fmt.Printf("[Speaker %d] %s\n", segment.Speaker, segment.Text)
}

// Custom SRT with speakers
func generateSRTWithSpeakers(result *omnivoice.TranscriptionResult) string {
    var builder strings.Builder
    for i, segment := range result.Segments {
        builder.WriteString(fmt.Sprintf("%d\n", i+1))
        builder.WriteString(fmt.Sprintf("%s --> %s\n",
            formatTime(segment.Start), formatTime(segment.End)))
        builder.WriteString(fmt.Sprintf("[Speaker %d]: %s\n\n",
            segment.Speaker, segment.Text))
    }
    return builder.String()
}
```

## Styling VTT

VTT supports styling with CSS:

```go
vtt := `WEBVTT

STYLE
::cue {
  background-color: rgba(0, 0, 0, 0.8);
  color: white;
  font-family: Arial, sans-serif;
}

::cue(.speaker1) {
  color: #00ff00;
}

::cue(.speaker2) {
  color: #ff0000;
}

00:00:00.000 --> 00:00:03.500
<c.speaker1>Hello, welcome to the show.</c>

00:00:03.500 --> 00:00:07.200
<c.speaker2>Thanks for having me!</c>
`
```

## Real-Time Subtitles

For live streaming, generate subtitles from streaming STT:

```go
stream, _ := provider.TranscribeStream(ctx, config)

var currentCue strings.Builder
var cueStart time.Duration

for result := range stream.Results() {
    if result.IsFinal {
        // Output completed cue
        fmt.Printf("%s --> %s\n%s\n\n",
            formatDuration(cueStart),
            formatDuration(result.End),
            currentCue.String())

        currentCue.Reset()
        cueStart = result.End
    } else {
        // Update interim text
        currentCue.WriteString(result.Text)
    }
}
```

## Best Practices

1. **Always enable word timestamps** - Required for accurate timing
2. **Use smart formatting** - Proper punctuation improves readability
3. **Test line lengths** - Different players have different limits
4. **Consider reading speed** - ~150-180 words per minute is comfortable
5. **Handle edge cases** - Empty segments, very short words

## Common Issues

### No Word Timestamps

```go
if len(result.Words) == 0 {
    // Fall back to segment-based timing
    for _, segment := range result.Segments {
        // Use segment.Start and segment.End
    }
}
```

### Overlapping Cues

```go
// Ensure minimum gap between cues
const minGap = 100 * time.Millisecond

prevEnd := time.Duration(0)
for _, cue := range cues {
    if cue.Start < prevEnd+minGap {
        cue.Start = prevEnd + minGap
    }
    prevEnd = cue.End
}
```

## Next Steps

- [STT Guide](stt.md) - Transcription details
- [Streaming](streaming.md) - Real-time subtitles
- [Voice Agents](voice-agents.md) - Live captioning for calls
