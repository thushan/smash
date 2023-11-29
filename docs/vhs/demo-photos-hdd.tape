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
Type "./smash /media/thushan/smash/photos/ --exclude-dir=_uploaded,sort,tmp,events"
Sleep 500ms
Enter
Sleep 30s
