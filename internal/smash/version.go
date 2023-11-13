package smash

import (
	"log"

	"github.com/thushan/smash/internal/theme"
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
	log.Println(theme.ColourSplash(`╔───────────────────────────────────────────────╗
│  ███████╗███╗   ███╗ █████╗ ███████╗██╗  ██╗  │
│  ██╔════╝████╗ ████║██╔══██╗██╔════╝██║  ██║  │
│  ███████╗██╔████╔██║███████║███████╗███████║  │
│  ╚════██║██║╚██╔╝██║██╔══██║╚════██║██╔══██║  │
│  ███████║██║ ╚═╝ ██║██║  ██║███████║██║  ██║  │
│  ╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝  │`))
	log.Println(theme.ColourSplash("│ "), theme.StyleUrl(Home), "           ", theme.ColourVersion(Version), theme.ColourSplash(" │"))
	log.Println(theme.ColourSplash(`╚───────────────────────────────────────────────╝`))

	if extendedInfo {
		log.Println("Commit: ", theme.ColourVersionMeta(Commit))
		log.Println("Built:  ", theme.ColourVersionMeta(Date))
		log.Println("Using:  ", theme.ColourVersionMeta(User))
	}
}
