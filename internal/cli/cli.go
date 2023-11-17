package cli

import (
	"errors"
	"log"
	"os"
	"runtime"

	"github.com/thushan/smash/internal/theme"

	"github.com/thushan/smash/internal/algorithms"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
	"github.com/thushan/smash/internal/smash"
)

var (
	af      *smash.Flags
	rootCmd = &cobra.Command{
		Use:          "smash [flags] [locations-to-smash]",
		Short:        "Find duplicates fast!",
		Long:         "",
		Version:      smash.Version,
		SilenceUsage: true,
		RunE:         runE,
	}
)

func init() {
	af = &smash.Flags{}
	rootCmd.SilenceErrors = true
	rootCmd.PersistentFlags().Var(
		enumflag.New(&af.Algorithm, "algorithm", algorithms.HashAlgorithms, enumflag.EnumCaseInsensitive),
		"algorithm",
		"Algorithm to use, can be 'xxhash', 'fnv128', 'fnv128a'")
	flags := rootCmd.Flags()
	flags.StringSliceVarP(&af.Base, "base", "", nil, "Base directories to use for comparison. Eg. --base=/c/dos,/c/run/dos/")
	flags.StringSliceVarP(&af.ExcludeFile, "exclude-file", "", nil, "Files to exclude separated by comma. Eg. --exclude-file=.gitignore,*.csv")
	flags.StringSliceVarP(&af.ExcludeDir, "exclude-dir", "", nil, "Directories to exclude separated by comma. Eg. --exclude-dir=.git,.idea")
	flags.IntVarP(&af.MaxThreads, "max-threads", "p", runtime.NumCPU(), "Maximum threads to utilise.")
	flags.IntVarP(&af.MaxWorkers, "max-workers", "w", bestMaxWorkers(), "Maximum workers to utilise when smashing.")
	flags.BoolVarP(&af.DisableSlicing, "disable-slicing", "", false, "Disable slicing (hashes full file).")
	flags.BoolVarP(&af.IgnoreEmptyFiles, "ignore-emptyfiles", "", false, "Ignore & don't report on empty/zero byte files.")
	flags.BoolVarP(&af.Silent, "silent", "q", false, "Run in silent mode.")
	flags.BoolVarP(&af.Verbose, "verbose", "", false, "Run in verbose mode.")
}
func bestMaxWorkers() int {
	cpus := runtime.NumCPU()
	if cpus < 6 {
		return 2
	} else {
		return cpus / 2
	}
}
func Main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(os.Stdout)
	if err := rootCmd.Execute(); err != nil {
		theme.Error.Println(err)
		os.Exit(1)
	}
}

func runE(command *cobra.Command, args []string) error {
	var locations []string
	if len(args) == 0 {
		// If no path found take the current path
		if wd, err := os.Getwd(); err != nil {
			locations = []string{"."}
		} else {
			locations = []string{wd}
		}
	} else {
		locations = verifyLocations(append(args, af.Base...), af.Silent)
	}

	if len(locations) == 0 {
		return errors.New("No valid locations to smash :(")
	}

	a := smash.App{
		Flags:     af,
		Args:      args,
		Locations: locations,
	}
	return a.Run()
}

func verifyLocations(locations []string, silent bool) []string {
	vl := locations[:0]
	for _, location := range locations {
		if _, err := os.Stat(location); os.IsNotExist(err) {
			if !silent {
				theme.Warn.Println("Ignoring invalid path ", theme.ColourFilename(location))
			}
			continue
		}
		vl = append(vl, location)
	}
	return vl
}
