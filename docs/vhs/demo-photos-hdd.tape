###
# VHS Tape for recording smash
# https://github.com/charmbracelet/vhs
###

Output docs/artefacts/smash-v0.0.x-hdd-photos-demo.gif

# Setup
Set FontSize 24
Set Width 1500
Set Height 1200
Set WindowBar Colorful
Set Theme "TokyoNight"

# smash Linux/drivers
Type "./smash /media/thushan/smash/photos/ --exclude-dir=sort,tmp,events -o report.json"
Sleep 500ms
Enter
Sleep 30s
Type "clear"
Sleep 500ms
Enter
Type `jq '.analysis.dupes[]|[.location,.path,.filename]|join("/")' report.json`
Sleep 500ms
Enter
Sleep 5s
Type "clear"
Sleep 500ms
Enter
Type `jq '.analysis.dupes[]|[.location,.path,.filename]|join("/")' report.json | xargs ls -lh`
Sleep 500ms
Enter
Sleep 5s

