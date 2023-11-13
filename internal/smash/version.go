package smash

import (
	"github.com/thushan/smash/internal/theme/colour"
	"log"
)

var (
	Version string = "v0.0.2"
	Commit  string = "none"
	Date    string = "unknown"
	Home    string = "github.com/thushan/smash"
	Time    string = "nowish"
	User    string = "local"
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

	if extendedInfo {
		log.Println("Commit: ", colour.VersionMeta(Commit))
		log.Println("Built:  ", colour.VersionMeta(Date))
		log.Println("Using:  ", colour.VersionMeta(User))
	}
}
