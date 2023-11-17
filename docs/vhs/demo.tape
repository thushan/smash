###
# VHS Tape for recording smash
# https://github.com/charmbracelet/vhs
###

Output docs/artefacts/smash-v0.0.x-demo.gif

# Setup
Set FontSize 20
Set Width 1500
Set Height 1000
Set WindowBar Colorful
Set Theme "JetBrains Darcula"

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
Type "./smash ~/linux/drivers
Sleep 500ms
Enter
Sleep 10s
