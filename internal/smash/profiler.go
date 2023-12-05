package smash

import (
	"log"
	"net/http"
	"net/http/pprof"
)

func InitialiseProfiler() {
	http.DefaultServeMux = http.NewServeMux()
	go func() {
		http.HandleFunc("/debug/pprof/", pprof.Index)
		http.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		http.HandleFunc("/debug/pprof/profile", pprof.Profile)
		http.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		http.HandleFunc("/debug/pprof/trace", pprof.Trace)
		log.Fatal(http.ListenAndServe("localhost:1984", nil))
	}()
}
