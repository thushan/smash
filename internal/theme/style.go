package theme

import (
	"github.com/pterm/pterm"
)

var (
	WithContextGlyph = "└─"

	SequenceIndexing = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	SequenceSmashing = []string{"◰", "◳", "◲", "◱"}
	SequenceFinalise = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	SequenceInternet = []string{"🌍", "🌎", "🌏"}
	SequenceTimeSoon = []string{"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛"}
	SequenceTimeLong = []string{"🕐", "🕜", "🕑", "🕝", "🕒", "🕞", "🕓", "🕟", "🕔", "🕠", "🕕", "🕡", "🕖", "🕢", "🕗", "🕣", "🕘", "🕤", "🕙", "🕥", "🕚", "🕦", "🕛", "🕧"}

	SequenceSmashingAlt = []string{"⬒", "⬔", "⬓", "⬕"}

	Info  = pterm.Info
	Error = pterm.Error
	Warn  = pterm.Warning
	Fatal = pterm.Error.WithFatal(true)

	Verbose      *pterm.PrefixPrinter
	WarnSkipping *pterm.PrefixPrinter

	StyleBold    *pterm.Style
	StyleHeading *pterm.Style
	StyleContext *pterm.Style
)

func init() {

	skippingPrefix := pterm.Warning.Prefix
	skippingPrefix.Text = "SKIP"
	WarnSkipping = pterm.Warning.WithPrefix(skippingPrefix)

	verbosePrefix := pterm.Info.Prefix
	verbosePrefix.Text = "VERBOSE"
	Verbose = pterm.Info.WithPrefix(verbosePrefix)

	StyleBold = pterm.NewStyle(pterm.Bold)
	StyleHeading = pterm.NewStyle(pterm.FgCyan, pterm.Bold)
	StyleContext = pterm.NewStyle(pterm.FgDarkGray, pterm.Italic)

}

func Println(message ...any) {
	pterm.Println(message...)
}
func PrintlnWithContext(context string, message ...any) {
	pterm.Println(pterm.Sprintln(context), StyleContext.Sprint(WithContextGlyph), StyleContext.Sprint(message...))
}
func WarnSkipWithContext(context string, message ...any) {
	pterm.Println(WarnSkipping.Sprintln(context), StyleContext.Sprint(WithContextGlyph), StyleContext.Sprint(message...))
}

func MultiWriter() pterm.MultiPrinter {
	return pterm.DefaultMultiPrinter
}
func DefaultSpinner() pterm.SpinnerPrinter {
	spinner := pterm.DefaultSpinner
	return spinner
}
func IndexingSpinner() pterm.SpinnerPrinter {
	spinner := DefaultSpinner()
	spinner.Sequence = SequenceIndexing
	return spinner
}
func SmashingSpinner() pterm.SpinnerPrinter {
	spinner := DefaultSpinner()
	spinner.Sequence = SequenceSmashing
	return spinner
}
func FinaliseSpinner() pterm.SpinnerPrinter {
	spinner := DefaultSpinner()
	spinner.Sequence = SequenceFinalise
	return spinner
}
func InternetSpinner() pterm.SpinnerPrinter {
	spinner := DefaultSpinner()
	spinner.Sequence = SequenceInternet
	return spinner
}
func TimeSoonSpinner() pterm.SpinnerPrinter {
	spinner := DefaultSpinner()
	spinner.Sequence = SequenceTimeSoon
	return spinner
}
func TimeLongSpinner() pterm.SpinnerPrinter {
	spinner := DefaultSpinner()
	spinner.Sequence = SequenceTimeLong
	return spinner
}
func ColourError(message ...any) string {
	return pterm.Red(message...)
}
func ColourSplash(message ...any) string {
	return pterm.LightGreen(message...)
}
func ColourPath(message ...any) string {
	return pterm.Blue(message...)
}

func StyleUrl(message ...any) string {
	return pterm.LightBlue(message...)
}
func Hyperlink(uri string, text string) string {
	return "\x1b]8;;" + uri + "\x07" + text + "\x1b]8;;\x07" + "\u001b[0m"
}
func ColourFilename(message ...any) string {
	return pterm.LightMagenta(message...)
}
func ColourFilenameA(message ...any) string {
	return pterm.Magenta(message...)
}
func ColourFileSize(message ...any) string {
	return pterm.Blue(message...)
}
func ColourFileSizeA(message ...any) string {
	return pterm.Cyan(message...)
}
func ColourHash(message ...any) string {
	return pterm.Gray(message...)
}
func ColourVersion(message ...any) string {
	return pterm.LightYellow(message...)
}
func ColourVersionMeta(message ...any) string {
	return pterm.Magenta(message...)
}
func ColourFolderHierarchy(message ...any) string {
	return pterm.Yellow(message...)
}
func ColourSuccess(message ...any) string {
	return pterm.Green(message...)
}
func ColourTime(message ...any) string {
	return pterm.Green(message...)
}
func ColourNumber(message ...any) string {
	return pterm.Blue(message...)
}
func ColourConfig(message ...any) string {
	return pterm.Magenta(message...)
}
func ColourConfigA(message ...any) string {
	return pterm.LightYellow(message...)
}
