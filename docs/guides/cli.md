# Command-Line Interface

OmniVoice includes a CLI for speech-to-text transcription and text-to-speech
synthesis without writing code.

## Installation

```bash
go install github.com/plexusone/omnivoice/cmd/omnivoice@latest
```

## Commands

### transcribe

Transcribe audio files to text.

```bash
omnivoice transcribe [flags] <audio-file>
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--provider` | `-p` | `deepgram` | STT provider (deepgram, openai, elevenlabs) |
| `--output` | `-o` | stdout | Output file path |
| `--format` | `-f` | `text` | Output format (text, json, srt, vtt) |
| `--language` | `-l` | `en-US` | Language code (BCP-47) |
| `--model` | `-m` | | Provider-specific model |
| `--diarize` | | `false` | Enable speaker diarization |
| `--timestamps` | | `false` | Enable word timestamps |
| `--verbose` | `-v` | `false` | Verbose output |
| `--quiet` | `-q` | `false` | Suppress non-essential output |

### speak

Synthesize speech from text and write it to an audio file.

```bash
omnivoice speak [flags] [text]
```

The text comes from the argument, or from standard input with `-` or no
argument. The output's extension picks the format: `.mp3`, `.wav`, `.opus` or
`.pcm`. Raw PCM from a provider is written to a `.wav` with a WAV header. A
provider that can't return the format asked for, such as ElevenLabs for
`.opus`, stops with an error and writes nothing.

Set the provider's API key in its environment variable, such as
`ELEVENLABS_API_KEY`. Providers that need none, such as `twilio`, work
without it.

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--voice` | | | Voice ID (required; list them with `omnivoice voices`) |
| `--output` | `-o` | | Audio file to write (required) |
| `--provider` | `-p` | `elevenlabs` | TTS provider |
| `--model` | `-m` | | Provider-specific model |
| `--locale` | `-l` | `en-US` | BCP-47 locale of the text, for choosing pronunciations |
| `--terms` | | | A terminology-spec terms directory to take pronunciations from |
| `--ipa` | | | `TERM=IPA`, such as `Miura=ˈmjura`; repeatable |
| `--say` | | | `TERM=RESPELLING`, such as `Miura=mee oora`; repeatable |
| `--phonemes` | | `auto` | ElevenLabs: send IPA as phoneme tags `auto` (by model), `on` or `off` |

### voices

List a TTS provider's voices.

```bash
omnivoice voices [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--provider` | `-p` | `elevenlabs` | TTS provider |
| `--match` | | | Only voices whose name, language or details contain this text |
| `--json` | | `false` | Print JSON |

### providers list

List available STT providers and their configuration status.

```bash
omnivoice providers list
```

## Examples

### Basic Transcription

```bash
export DEEPGRAM_API_KEY="your-api-key"

# Output to terminal
omnivoice transcribe podcast.mp3

# Save to file
omnivoice transcribe -o transcript.txt podcast.mp3
```

### JSON Output (OmniVoice Transcript Format)

Get full metadata including timestamps, speakers, and confidence scores:

```bash
omnivoice transcribe -p deepgram --diarize --timestamps -f json -o transcript.json podcast.mp3
```

Output:

```json
{
  "$schema": "https://omnivoice.dev/schema/transcript-v1.json",
  "version": "1.0",
  "text": "Hello and welcome to the podcast...",
  "language": "en-US",
  "duration_ms": 330000,
  "segments": [
    {
      "text": "Hello and welcome to the podcast.",
      "start_ms": 0,
      "end_ms": 2500,
      "speaker": "Speaker 1",
      "words": [
        {"text": "Hello", "start_ms": 0, "end_ms": 300}
      ]
    }
  ],
  "metadata": {
    "provider": "deepgram",
    "created_at": "2026-04-28T10:30:00Z",
    "audio_file": "podcast.mp3"
  }
}
```

### Subtitle Generation

#### SRT (SubRip)

```bash
omnivoice transcribe -p deepgram -f srt -o subtitles.srt podcast.mp3
```

Output:

```
1
00:00:00,000 --> 00:00:02,500
Hello and welcome to the podcast.

2
00:00:02,800 --> 00:00:05,200
Today we're talking about AI.
```

#### WebVTT

```bash
omnivoice transcribe -p deepgram -f vtt -o subtitles.vtt podcast.mp3
```

Output:

```
WEBVTT

1
00:00:00.000 --> 00:00:02.500
Hello and welcome to the podcast.

2
00:00:02.800 --> 00:00:05.200
Today we're talking about AI.
```

### Speaker Diarization

Identify different speakers in the audio:

```bash
omnivoice transcribe -p deepgram --diarize -f json -o meeting.json meeting.mp3
```

### Using Different Providers

```bash
# Deepgram (default)
export DEEPGRAM_API_KEY="your-key"
omnivoice transcribe -p deepgram audio.mp3

# OpenAI Whisper
export OPENAI_API_KEY="your-key"
omnivoice transcribe -p openai audio.mp3

# ElevenLabs
export ELEVENLABS_API_KEY="your-key"
omnivoice transcribe -p elevenlabs audio.mp3
```

### Speech Synthesis

```bash
# Find a voice with an Italian accent
omnivoice voices --match italian

# Synthesize a sentence
omnivoice speak --voice <id> -o miura.mp3 "The Miura was shown at Geneva."

# Narration from a file, with pronunciations from a terms directory
omnivoice speak --voice <id> --terms terms -o lesson.mp3 - < lesson.txt
```

### Pronouncing Names

Give a term's pronunciation as IPA, a respelling, or both. ElevenLabs uses
the IPA where the model accepts phonemes and the respelling otherwise; see
[Pronunciations](tts.md#pronunciations). For other providers, `speak` puts
the respellings into the text itself, and a term with only IPA is spoken as
written.

```bash
omnivoice speak --voice <id> --model eleven_flash_v2 \
  --ipa "Lamborghini=lamborˈɡiːni" --say "Lamborghini=lambor gheeny" \
  -o lambo.mp3 "Tractors made Lamborghini wealthy."
```

`--ipa` and `--say` win over a `--terms` directory for the terms they name.
Slashes or brackets around IPA, as dictionaries print it, are removed.

### Check Provider Status

```bash
omnivoice providers list
```

Output:

```
PROVIDER    ENV VAR             CONFIGURED  FEATURES
--------    -------             ----------  --------
deepgram    DEEPGRAM_API_KEY    Yes         streaming, diarization, timestamps, punctuation
elevenlabs  ELEVENLABS_API_KEY  No          diarization, timestamps
openai      OPENAI_API_KEY      Yes         timestamps, punctuation
```

## Output Formats

| Format | Use Case | Features |
|--------|----------|----------|
| `text` | Simple transcript | Plain text only |
| `json` | Programmatic access | Full metadata, timestamps, speakers, confidence |
| `srt` | Video subtitles | Timed text, widely supported |
| `vtt` | Web video | W3C standard, HTML5 video support |

## Supported Audio Formats

- MP3
- WAV
- FLAC
- OGG
- M4A
- WebM

## Environment Variables

| Variable | Provider | Required |
|----------|----------|----------|
| `DEEPGRAM_API_KEY` | Deepgram | For Deepgram provider |
| `OPENAI_API_KEY` | OpenAI | For OpenAI provider |
| `ELEVENLABS_API_KEY` | ElevenLabs | For ElevenLabs provider |
