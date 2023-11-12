package theme

import "github.com/pterm/pterm"

var (
	SequenceIndexing = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	SequenceSmashing = []string{"◰", "◳", "◲", "◱"}
	SequenceFinalise = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	SequenceInternet = []string{"🌍", "🌎", "🌏"}
	SequenceTimeSoon = []string{"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛"}
	SequenceTimeLong = []string{"🕐", "🕜", "🕑", "🕝", "🕒", "🕞", "🕓", "🕟", "🕔", "🕠", "🕕", "🕡", "🕖", "🕢", "🕗", "🕣", "🕘", "🕤", "🕙", "🕥", "🕚", "🕦", "🕛", "🕧"}

	SequenceSmashingAlt = []string{"⬒", "⬔", "⬓", "⬕"}
)

func MultiWriter() pterm.MultiPrinter {
	return pterm.DefaultMultiPrinter
}
func DefaultSpinner() pterm.SpinnerPrinter {
	return pterm.DefaultSpinner
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
