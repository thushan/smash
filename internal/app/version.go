package app

import (
	"log"

	"github.com/logrusorgru/aurora/v3"
)

var (
	Version = "v1.0.0"
	Edition = "open-source"
	Home    = "github.com/thushan/smash"
	Time    string
	User    string
)

func PrintVersionInfo(extendedInfo bool) {
	log.Println(aurora.Green(`        
╔───────────────────────────────────────────────╗
│  ███████╗███╗   ███╗ █████╗ ███████╗██╗  ██╗  │
│  ██╔════╝████╗ ████║██╔══██╗██╔════╝██║  ██║  │
│  ███████╗██╔████╔██║███████║███████╗███████║  │
│  ╚════██║██║╚██╔╝██║██╔══██║╚════██║██╔══██║  │
│  ███████║██║ ╚═╝ ██║██║  ██║███████║██║  ██║  │
│  ╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝  │`))
	log.Println(aurora.Green("│ "), aurora.Yellow(Home), "           ", aurora.Blue(Version), aurora.Green(" │"))
	log.Println(aurora.Green(`╚───────────────────────────────────────────────╝`))

	if extendedInfo {
		log.Println("Edition: ", Edition)
		log.Println("by: ", User)
		log.Println("on: ", Time)
	}
}
