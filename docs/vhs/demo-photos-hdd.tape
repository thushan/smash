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
Type "# In this demo, we're using a slow USB SATA 2.5inch disk with photos"
Enter
Type "# but we want to ignore sort, tmp & events folders in the collection"
Enter
Sleep 1s
Type "./smash /media/thushan/smash/photos/ --exclude-dir=sort,tmp,events -o report.json"
Sleep 500ms
Enter
Sleep 30s
Type "clear"
Sleep 500ms
Enter
Type "# So let's see the main files that were duplicated:"
Enter
Sleep 1s
Type `jq '.analysis.dupes[]|[.location,.path,.filename]|join("/")' report.json`
Sleep 500ms
Enter
Sleep 2s
Type "# and now the duplicates:"
Enter
Sleep 1s
Type `jq '.analysis.dupes[].duplicates[]|[.location,.path,.filename]|join("/")' report.json`
Sleep 500ms
Enter
Sleep 5s
Type "clear"
Sleep 500ms
Enter
Type "# You could now remove the files but let's just list the filesizes for now:"
Enter
Sleep 1s
Type `jq '.analysis.dupes[].duplicates[]|[.location,.path,.filename]|join("/")' report.json | xargs exa -lh --no-user --no-permissions`
Sleep 500ms
Enter
Sleep 5s

