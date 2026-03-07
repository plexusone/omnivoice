// Package omnivoice provides a unified interface for speech-to-text and text-to-speech.
// This is the batteries-included package that imports all providers.
// For a minimal dependency footprint, use github.com/plexusone/omnivoice-core instead.
package omnivoice

import (
	"github.com/plexusone/omnivoice-core/callsystem"
)

// Re-export CallSystem types from omnivoice-core
type (
	// Call represents a phone or video call.
	Call = callsystem.Call

	// CallSystem defines the interface for telephony/meeting integrations.
	CallSystem = callsystem.CallSystem

	// CallStatus represents the call state.
	CallStatus = callsystem.CallStatus

	// CallDirection indicates inbound or outbound call.
	CallDirection = callsystem.CallDirection

	// CallHandler is called when a new call arrives.
	CallHandler = callsystem.CallHandler

	// CallSystemConfig configures a call system integration.
	CallSystemConfig = callsystem.CallSystemConfig

	// CallOption configures an outbound call.
	CallOption = callsystem.CallOption

	// CallOptions holds parsed options for MakeCall.
	CallOptions = callsystem.CallOptions
)

// Re-export CallStatus constants
const (
	StatusRinging  = callsystem.StatusRinging
	StatusAnswered = callsystem.StatusAnswered
	StatusEnded    = callsystem.StatusEnded
	StatusFailed   = callsystem.StatusFailed
	StatusBusy     = callsystem.StatusBusy
	StatusNoAnswer = callsystem.StatusNoAnswer
	CallInbound    = callsystem.Inbound
	CallOutbound   = callsystem.Outbound
)

// Re-export CallOption functions
var (
	WithFrom             = callsystem.WithFrom
	WithTimeout          = callsystem.WithTimeout
	WithMachineDetection = callsystem.WithMachineDetection
	WithRecording        = callsystem.WithRecording
	WithWhisper          = callsystem.WithWhisper
	WithAgent            = callsystem.WithAgent
	WithStatusCallback   = callsystem.WithStatusCallback
)
