package nerdstats

import (
	"runtime"
)

type NerdStats struct {
	Allocations,
	TotalAllocations,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	GcPauseTotalNs uint64

	CompletedGcCycles uint32
	GoRoutines        int
}

func Snapshot() NerdStats {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	return NerdStats{
		GoRoutines:        runtime.NumGoroutine(),
		Allocations:       rtm.Alloc,
		TotalAllocations:  rtm.TotalAlloc,
		Sys:               rtm.Sys,
		Mallocs:           rtm.Mallocs,
		Frees:             rtm.Frees,
		LiveObjects:       rtm.Mallocs - rtm.Frees,
		GcPauseTotalNs:    rtm.PauseTotalNs,
		CompletedGcCycles: rtm.NumGC,
	}
}
