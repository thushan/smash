package smash

import (
	"fmt"
	"log"

	"github.com/thushan/smash/internal/theme"
)

var (
	Version = "v0.7.0"
	Commit  = "none"
	Date    = "unknown"
	Time    = "nowish"
	User    = "local"
)

const (
	GithubHomeText  = "github.com/thushan/smash"
	GithubHomeUri   = "https://github.com/thushan/smash"
	GithubLatestUri = "https://github.com/thushan/smash/releases/latest"
)

func PrintVersionInfo(extendedInfo bool) {
	githubUri := theme.Hyperlink(GithubHomeUri, GithubHomeText)
	latestUri := theme.Hyperlink(GithubLatestUri, Version)
	padLatest := fmt.Sprintf("%*s", 17-len(Version), "")

	log.Println(theme.ColourSplash(`╔───────────────────────────────────────────────╗
│  ███████╗███╗   ███╗ █████╗ ███████╗██╗  ██╗  │
│  ██╔════╝████╗ ████║██╔══██╗██╔════╝██║  ██║  │
│  ███████╗██╔████╔██║███████║███████╗███████║  │
│  ╚════██║██║╚██╔╝██║██╔══██║╚════██║██╔══██║  │
│  ███████║██║ ╚═╝ ██║██║  ██║███████║██║  ██║  │
│  ╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝  │`))
	log.Println(theme.ColourSplash("│ "), theme.StyleUrl(githubUri), padLatest, theme.ColourVersion(latestUri), theme.ColourSplash(" │"))
	log.Println(theme.ColourSplash(`╚───────────────────────────────────────────────╝`))
	if extendedInfo {
		log.Println(" Commit:", theme.ColourVersionMeta(Commit))
		log.Println("  Built:", theme.ColourVersionMeta(Date))
		log.Println("  Using:", theme.ColourVersionMeta(User))
	}
}
