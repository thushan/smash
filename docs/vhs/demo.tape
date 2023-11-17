###
# VHS Tape for recording smash
# https://github.com/charmbracelet/vhs
###

Output docs/artefacts/smash-v0.0.x-demo.gif

# Setup
Set FontSize 24
Set Width 1600
Set Height 900
Set WindowBar Colorful
Set Theme "TokyoNight"

# Version Check
Type "./smash --version"
Sleep 500ms
Enter
Sleep 2s

# Clear
Type "clear"
Sleep 500ms
Enter

# smash Linux/drivers
Type "./smash /linux/drivers
Sleep 500ms
Enter
Sleep 60s
