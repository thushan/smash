package cli

import (
	"github.com/spf13/cobra"
	"github.com/thushan/smash/internal/app"
	"github.com/thushan/smash/internal/smash"
	"log"
	"os"
)

var af *app.Flags
var rootCmd = &cobra.Command{
	Use:          "smash [flags] [locations-to-scan]",
	Short:        "Find duplicates fast!",
	Long:         "",
	Version:      smash.Version,
	SilenceUsage: true,
	RunE:         runE,
}

func init() {
	af = &app.Flags{}
	flags := rootCmd.Flags()
	flags.BoolVarP(&af.Silent, "silent", "q", false, "Run in silent mode.")
	flags.BoolVarP(&af.Verbose, "verbose", "", false, "Run in verbose mode.")
}

func Main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	smash.PrintVersionInfo(false)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
func runE(command *cobra.Command, args []string) error {
	var locations []string
	if len(args) == 0 {
		locations = []string{"."}
	}

	a := smash.App{
		Flags:     af,
		Args:      args,
		Locations: locations,
	}
	return a.Run()
}
