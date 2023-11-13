package colour

import "github.com/pterm/pterm"

func Heading(message ...any) string {
	return pterm.Cyan(message...)
}

func Error(message ...any) string {
	return pterm.Red(message...)
}

func Splash(message ...any) string {
	return pterm.LightGreen(message...)
}

func Url(message ...any) string {
	return pterm.LightBlue(message...)
}

func Version(message ...any) string {
	return pterm.LightYellow(message...)
}
func VersionMeta(message ...any) string {
	return pterm.Magenta(message...)
}
