package smash

import (
	"log"
	"net/http"
	"net/http/pprof"
	"time"
)

func InitialiseProfiler() {
	http.DefaultServeMux = http.NewServeMux()
	go func() {
		address := "localhost:1984"
		server := &http.Server{
			Addr:         address,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		http.HandleFunc("/debug/pprof/", pprof.Index)
		http.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		http.HandleFunc("/debug/pprof/profile", pprof.Profile)
		http.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		http.HandleFunc("/debug/pprof/trace", pprof.Trace)
		log.Fatal(server.ListenAndServe())
	}()
}
