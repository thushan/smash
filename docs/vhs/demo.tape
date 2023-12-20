###
# VHS Tape for recording smash
# https://github.com/charmbracelet/vhs
###

Output docs/artefacts/smash-v0.0.x-demo.gif

# Setup
Set FontSize 24
Set Width 1500
Set Height 1200
Set WindowBar Colorful
Set Theme "JetBrains Darcula"

# smash Linux/drivers
Type "./smash ~/linux/drivers --exclude-dir=git -o report.json"
Sleep 500ms
Enter
Sleep 10s
Type "clear"
Sleep 500ms
Enter
Type `jq '.analysis.dupes[]|[.location,.path,.filename]|join("/")' report.json | xargs wc -l`
Sleep 500ms
Enter
Sleep 5s
Type `jq '.analysis.dupes[]|[.location,.path,.filename]|join("/")' report.json | xargs rm`
Sleep 500ms
Enter
Sleep 5s
Type "cd ~/linux/drivers"
Sleep 500ms
Enter
Sleep 2s
Type "git status -s"
Sleep 500ms
Enter
Sleep 3s
Type "git reset --hard"
Sleep 500ms
Enter
Sleep 3s
Type "git status -s"
Sleep 500ms
Enter
Sleep 5s
