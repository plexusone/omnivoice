package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/plexusone/omnivoice"
	_ "github.com/plexusone/omnivoice/providers/all"
	"github.com/spf13/cobra"
)

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "Manage STT providers",
	Long:  `Commands for listing and managing speech-to-text providers.`,
}

var providersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available STT providers",
	Long: `List all available speech-to-text providers with their configuration status.

Shows provider name, required environment variable, and whether the API key is configured.`,
	RunE: runProvidersList,
}

func init() {
	providersCmd.AddCommand(providersListCmd)
}

// providerInfo holds metadata about a provider.
type providerInfo struct {
	Name       string
	EnvVar     string
	Features   []string
	Configured bool
}

// knownProviders returns metadata for known STT providers.
func knownProviders() []providerInfo {
	return []providerInfo{
		{
			Name:     "deepgram",
			EnvVar:   "DEEPGRAM_API_KEY",
			Features: []string{"streaming", "diarization", "timestamps", "punctuation"},
		},
		{
			Name:     "openai",
			EnvVar:   "OPENAI_API_KEY",
			Features: []string{"timestamps", "punctuation"},
		},
		{
			Name:     "elevenlabs",
			EnvVar:   "ELEVENLABS_API_KEY",
			Features: []string{"diarization", "timestamps"},
		},
	}
}

func runProvidersList(cmd *cobra.Command, args []string) error {
	// Get registered providers
	registered := omnivoice.ListSTTProviders()
	registeredSet := make(map[string]bool)
	for _, name := range registered {
		registeredSet[name] = true
	}

	// Build provider list with status
	providers := knownProviders()
	for i := range providers {
		providers[i].Configured = os.Getenv(providers[i].EnvVar) != ""
	}

	// Add any registered providers not in known list
	for _, name := range registered {
		found := false
		for _, p := range providers {
			if strings.EqualFold(p.Name, name) {
				found = true
				break
			}
		}
		if !found {
			envVar := strings.ToUpper(name) + "_API_KEY"
			providers = append(providers, providerInfo{
				Name:       name,
				EnvVar:     envVar,
				Features:   []string{},
				Configured: os.Getenv(envVar) != "",
			})
		}
	}

	// Sort by name
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Name < providers[j].Name
	})

	// Print table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PROVIDER\tENV VAR\tCONFIGURED\tFEATURES")
	fmt.Fprintln(w, "--------\t-------\t----------\t--------")

	for _, p := range providers {
		status := "No"
		if p.Configured {
			status = "Yes"
		}
		features := strings.Join(p.Features, ", ")
		if features == "" {
			features = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Name, p.EnvVar, status, features)
	}

	return w.Flush()
}
