// Package cli provides the command-line interface for omnivoice.
package cli

import (
	"fmt"
	"os"

	"github.com/plexusone/omnivoice"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	quiet   bool
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "omnivoice",
	Short: "Multi-provider speech-to-text CLI",
	Long: `omnivoice is a command-line tool for speech-to-text transcription.

It supports multiple providers including Deepgram, OpenAI, and ElevenLabs,
with features like speaker diarization, timestamps, and subtitle generation.`,
	Version: omnivoice.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential output")

	rootCmd.AddCommand(transcribeCmd)
	rootCmd.AddCommand(providersCmd)

	// Set version template
	rootCmd.SetVersionTemplate(fmt.Sprintf("omnivoice version %s\n", omnivoice.Version))
}

// logInfo prints informational messages unless quiet mode is enabled.
func logInfo(format string, args ...any) {
	if !quiet {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

// logVerbose prints verbose messages only when verbose mode is enabled.
func logVerbose(format string, args ...any) {
	if verbose && !quiet {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}
