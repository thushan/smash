package cli

import (
	"errors"
	"log"
	"os"
	"runtime"

	"github.com/thushan/smash/internal/algorithms"
	"github.com/thushan/smash/internal/smash"
	"github.com/thushan/smash/internal/theme"
	"github.com/thushan/smash/pkg/indexer"
	"github.com/thushan/smash/pkg/slicer"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

var (
	af      *smash.Flags
	rootCmd = &cobra.Command{
		Use:          "smash [flags] [locations-to-smash]",
		Short:        "Find duplicates fast!",
		Long:         "",
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
		"Algorithm to use to hash files. Supported: xxhash, murmur3, md5, sha512, sha256 (full list, see readme)")
	flags := rootCmd.Flags()
	flags.StringSliceVarP(&af.Base, "base", "", nil, "Base directories to use for comparison Eg. --base=/c/dos,/c/dos/run/,/run/dos/run")
	flags.StringSliceVarP(&af.ExcludeFile, "exclude-file", "", nil, "Files to exclude separated by comma Eg. --exclude-file=.gitignore,*.csv")
	flags.StringSliceVarP(&af.ExcludeDir, "exclude-dir", "", nil, "Directories to exclude separated by comma Eg. --exclude-dir=.git,.idea")
	flags.IntVarP(&af.MaxThreads, "max-threads", "p", runtime.NumCPU(), "Maximum threads to utilise")
	flags.IntVarP(&af.MaxWorkers, "max-workers", "w", runtime.NumCPU(), "Maximum workers to utilise when smashing")
	flags.Int64VarP(&af.MinSize, "min-size", "G", 0, "Minimum file size to consider for hashing (in bytes)")
	flags.Int64VarP(&af.MaxSize, "max-size", "L", 0, "Maximum file size to consider for hashing (in bytes)")
	flags.IntVarP(&af.ProgressUpdate, "progress-update", "", 5, "Update progress every x seconds")
	flags.IntVarP(&af.ShowTop, "show-top", "", 10, "Show the top x duplicates")
	flags.BoolVarP(&af.HideTopList, "no-top-list", "", false, "Hides top x duplicates list")
	flags.BoolVarP(&af.ShowDuplicates, "show-duplicates", "", false, "Show full list of duplicates")
	flags.BoolVarP(&af.DisableSlicing, "disable-slicing", "", false, "Disable slicing & hash the full file instead")
	flags.BoolVarP(&af.DisableMeta, "disable-meta", "", false, "Disable storing of meta-data to improve hashing mismatches")
	flags.BoolVarP(&af.DisableAutoText, "disable-autotext", "", false, "Disable detecting text-files to opt for a full hash for those")
	flags.BoolVarP(&af.IgnoreEmpty, "ignore-empty", "", true, "Ignore empty/zero byte files")
	flags.BoolVarP(&af.IgnoreHidden, "ignore-hidden", "", true, "Ignore hidden files & folders Eg. files/folders starting with '.'")
	flags.BoolVarP(&af.IgnoreSystem, "ignore-system", "", true, "Ignore system files & folders Eg. '$MFT', '.Trash'")
	flags.BoolVarP(&af.Silent, "silent", "q", false, "Run in silent mode")
	flags.BoolVarP(&af.Recurse, "recurse", "r", false, "Recursively search directories for files")
	flags.BoolVarP(&af.Verbose, "verbose", "", false, "Run in verbose mode")
	flags.BoolVarP(&af.Profile, "profile", "", false, "Enable Go Profiler - see localhost:1984/debug/pprof")
	flags.BoolVarP(&af.HideProgress, "no-progress", "", false, "Disable progress updates")
	flags.BoolVarP(&af.HideOutput, "no-output", "", false, "Disable report output")
	flags.BoolVarP(&af.ShowNerdStats, "nerd-stats", "", false, "Show nerd stats")
	flags.BoolVarP(&af.ShowVersion, "version", "v", false, "Show version information")
	flags.StringVarP(&af.OutputFile, "output-file", "o", "", "Export analysis as JSON (generated automatically like ./report-*.json)")
	flags.IntVarP(&af.Slices, "slices", "", slicer.DefaultSlices, "Number of Slices to use")
	flags.Int64VarP(&af.SliceSize, "slice-size", "", slicer.DefaultSliceSize, "Size of a Slice (in bytes)")
	flags.Int64VarP(&af.SliceThreshold, "slice-threshold", "", slicer.DefaultThreshold, "Threshold to use for slicing (in bytes) - if file is smaller than this, it won't be sliced")
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

	var locations []indexer.LocationFS

	if len(args) == 0 {
		// If no path found take the current path
		if wd, err := os.Getwd(); err != nil {
			locations = []indexer.LocationFS{{
				FS:   os.DirFS("."),
				Name: ".",
				Kind: indexer.Local,
			}}
		} else {
			locations = []indexer.LocationFS{{
				FS:   os.DirFS(wd),
				Name: wd,
				Kind: indexer.Local,
			}}
		}
	} else {
		locations = verifyLocations(append(args, af.Base...), af.Silent)
	}

	if len(locations) == 0 {
		return errors.New("no valid locations to smash :(")
	}

	a := smash.App{
		Flags:     af,
		Args:      args,
		Locations: locations,
	}
	return a.Run()
}

func verifyLocations(locations []string, silent bool) []indexer.LocationFS {
	var vl []indexer.LocationFS
	for _, location := range locations {
		if _, err := os.Stat(location); os.IsNotExist(err) {
			if !silent {
				theme.Warn.Println("Ignoring invalid path ", theme.ColourFilename(location))
			}
			continue
		}
		l := indexer.NewLocationFS(indexer.Local, location, os.DirFS(location))
		vl = append(vl, *l)
	}
	return vl
}
