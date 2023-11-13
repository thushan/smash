package smash

import (
	"log"

	"github.com/thushan/smash/internal/theme/colour"
)

var (
	Version = "v0.0.3"
	Commit  = "none"
	Date    = "unknown"
	Home    = "github.com/thushan/smash"
	Time    = "nowish"
	User    = "local"
)

func PrintVersionInfo(extendedInfo bool) {
	log.Println(colour.Splash(`╔───────────────────────────────────────────────╗
│  ███████╗███╗   ███╗ █████╗ ███████╗██╗  ██╗  │
│  ██╔════╝████╗ ████║██╔══██╗██╔════╝██║  ██║  │
│  ███████╗██╔████╔██║███████║███████╗███████║  │
│  ╚════██║██║╚██╔╝██║██╔══██║╚════██║██╔══██║  │
│  ███████║██║ ╚═╝ ██║██║  ██║███████║██║  ██║  │
│  ╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝  │`))
	log.Println(colour.Splash("│ "), colour.Url(Home), "           ", colour.Version(Version), colour.Splash(" │"))
	log.Println(colour.Splash(`╚───────────────────────────────────────────────╝`))

	if !extendedInfo {
		log.Println("Commit: ", colour.VersionMeta(Commit))
		log.Println("Built:  ", colour.VersionMeta(Date))
		log.Println("Using:  ", colour.VersionMeta(User))
	}
}
