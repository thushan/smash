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

func Path(message ...any) string {
	return pterm.Blue(message...)
}

func Url(message ...any) string {
	return pterm.LightBlue(message...)
}

func Meta(message ...any) string {
	return pterm.Magenta(message...)
}

func Version(message ...any) string {
	return pterm.LightYellow(message...)
}
func VersionMeta(message ...any) string {
	return pterm.Magenta(message...)
}
