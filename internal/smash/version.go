package smash

import (
	"log"

	"github.com/logrusorgru/aurora/v3"
)

var (
	Version string = "v0.0.1"
	Commit  string = "none"
	Date    string = "unknown"
	Home    string = "github.com/thushan/smash"
	Time    string = "nowish"
	User    string = "local"
)

func PrintVersionInfo(extendedInfo bool) {
	log.Println(aurora.Green(`╔───────────────────────────────────────────────╗
│  ███████╗███╗   ███╗ █████╗ ███████╗██╗  ██╗  │
│  ██╔════╝████╗ ████║██╔══██╗██╔════╝██║  ██║  │
│  ███████╗██╔████╔██║███████║███████╗███████║  │
│  ╚════██║██║╚██╔╝██║██╔══██║╚════██║██╔══██║  │
│  ███████║██║ ╚═╝ ██║██║  ██║███████║██║  ██║  │
│  ╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝  │`))
	log.Println(aurora.Green("│ "), aurora.Yellow(Home), "           ", aurora.Blue(Version), aurora.Green(" │"))
	log.Println(aurora.Green(`╚───────────────────────────────────────────────╝`))

	if extendedInfo {
		log.Println("Commit: ", aurora.BrightBlack(Commit))
		log.Println("Built:  ", aurora.BrightBlack(Date))
		log.Println("Using:  ", aurora.BrightBlack(User))
		log.Println("")
	}
}
