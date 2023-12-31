package smash

import (
	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/internal/theme"
	"github.com/thushan/smash/pkg/nerdstats"
)

func PrintNerdStats(stats nerdstats.NerdStats, context string) {
	theme.StyleContext.Println(context)

	bold := theme.StyleBold
	theme.Println(bold.Sprint("             OSMemory:"), theme.ColourConfig(humanize.Bytes(stats.Sys), " (", stats.Sys, " Bytes)"))
	theme.Println(bold.Sprint("          Allocations:"), theme.ColourConfig(humanize.Bytes(stats.Allocations), " (", stats.Allocations, " Bytes)"))
	theme.Println(bold.Sprint("  Allocations (total):"), theme.ColourConfig(humanize.Bytes(stats.TotalAllocations), " (", stats.TotalAllocations, " Bytes)"))
	theme.Println(bold.Sprint("              mallocs:"), theme.ColourConfig(stats.Mallocs))
	theme.Println(bold.Sprint("                frees:"), theme.ColourConfig(stats.Frees))
	theme.Println(bold.Sprint("          LiveObjects:"), theme.ColourConfig(stats.LiveObjects))
	theme.Println(bold.Sprint("       GC-Pauses (ns):"), theme.ColourConfig(stats.GcPauseTotalNs))
	theme.Println(bold.Sprint("GC-Cycles (completed):"), theme.ColourConfig(stats.CompletedGcCycles))
	theme.Println(bold.Sprint("  GoRoutines (active):"), theme.ColourConfig(stats.GoRoutines))
}
