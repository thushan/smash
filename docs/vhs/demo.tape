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
Type "# Let's start by de-duping the linux source tree,"
Enter
Type "# exlude the git dir and saving to report.json!"
Enter
Sleep 1s
Type "./smash ~/linux/drivers -r --exclude-dir=git -o report.json"
Sleep 500ms
Enter
Sleep 10s
Type "clear"
Sleep 500ms
Enter
Type "# Now let's look at only the duplicate files"
Enter
Sleep 1s
Type `jq '.analysis.dupes[].duplicates[]|[.location,.path,.filename]|join("/")' report.json | xargs wc -l`
Sleep 500ms
Enter
Sleep 5s
Type "# *gulp* let's remove them duplicates, shall well?"
Enter
Sleep 1s
Type `jq '.analysis.dupes[].duplicates[]|[.location,.path,.filename]|join("/")' report.json | xargs rm`
Sleep 500ms
Enter
Sleep 5s
Type "# Okay, let's just double check they got removed"
Enter
Sleep 1s
Type "cd ~/linux/drivers"
Sleep 500ms
Enter
Sleep 1s
Type "git status -s"
Sleep 500ms
Enter
Sleep 3s
Type "# Alright, let's reset this for gits & giggles..."
Enter
Sleep 1s
Type "git reset --hard"
Sleep 500ms
Enter
Sleep 3s
Type "git status -s"
Sleep 500ms
Enter
Sleep 5s
