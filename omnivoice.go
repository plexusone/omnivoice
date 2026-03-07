// Package omnivoice provides a unified interface for speech-to-text and text-to-speech.
//
// This is the batteries-included package that imports all providers.
// For a minimal dependency footprint, use github.com/plexusone/omnivoice-core instead.
//
// # Quick Start
//
// Import the package with all providers:
//
//	import (
//	    "github.com/plexusone/omnivoice"
//	    _ "github.com/plexusone/omnivoice/providers/all"
//	)
//
// Or import specific providers:
//
//	import (
//	    "github.com/plexusone/omnivoice"
//	    openai "github.com/plexusone/omnivoice-openai/omnivoice"
//	)
//
// # Creating Providers
//
//	// OpenAI provider
//	sttProvider := openai.NewSTTProvider(apiKey)
//	ttsProvider := openai.NewTTSProvider(apiKey)
//
//	// Create multi-provider client
//	sttClient := omnivoice.NewSTTClient(sttProvider)
//	ttsClient := omnivoice.NewTTSClient(ttsProvider)
package omnivoice

// Version is the current version of the omnivoice package.
const Version = "0.6.0"
