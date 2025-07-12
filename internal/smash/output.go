package smash

import (
	"github.com/pterm/pterm"
)

// OutputConfig contains the resolved output configuration
type OutputConfig struct {
	Silent         bool
	ShowProgress   bool
	UseColor       bool
	Verbose        bool
	GenerateReport bool
}

// OutputManager handles all output operations respecting environment and flags
type OutputManager struct {
	env    EnvironmentInfo
	config OutputConfig
}

// NewOutputManager creates a new output manager based on flags and environment
func NewOutputManager(flags *Flags) *OutputManager {
	env := DetectEnvironment()

	config := OutputConfig{
		// Silent overrides everything
		Silent: flags.Silent,

		// Show progress only if allowed by environment and flags
		ShowProgress: !flags.HideProgress &&
			!flags.Silent &&
			env.ShouldAnimate,

		// Colors depend on environment and force flags
		UseColor: !env.NoColor &&
			(env.IsTerminal || env.ForceColor) &&
			!flags.Silent,

		// Verbose only works if not silent
		Verbose: flags.Verbose && !flags.Silent,

		// Report generation
		GenerateReport: !flags.HideOutput,
	}

	// Configure pterm based on environment
	if !config.UseColor {
		pterm.DisableColor()
	}

	if !env.IsTerminal || env.IsCI || env.IsTesting {
		pterm.DisableStyling()
	}

	// Disable all pterm output in test environments
	if env.IsTesting {
		pterm.DisableOutput()
	}

	return &OutputManager{
		config: config,
		env:    env,
	}
}

// SpinnerHandle provides a common interface for spinners
type SpinnerHandle interface {
	UpdateText(text string)
	Success(text string)
	Fail(text string)
	Warning(text string)
	Stop() error
}

// NoOpSpinner is a spinner that does nothing (for silent/non-TTY modes)
type NoOpSpinner struct{}

func (n *NoOpSpinner) UpdateText(text string) {}
func (n *NoOpSpinner) Success(text string)    {}
func (n *NoOpSpinner) Fail(text string)       {}
func (n *NoOpSpinner) Warning(text string)    {}
func (n *NoOpSpinner) Stop() error            { return nil }

// ptermSpinnerAdapter wraps pterm.SpinnerPrinter to implement SpinnerHandle
type ptermSpinnerAdapter struct {
	spinner *pterm.SpinnerPrinter
}

func (p *ptermSpinnerAdapter) UpdateText(text string) {
	p.spinner.UpdateText(text)
}

func (p *ptermSpinnerAdapter) Success(text string) {
	p.spinner.Success(text)
}

func (p *ptermSpinnerAdapter) Fail(text string) {
	p.spinner.Fail(text)
}

func (p *ptermSpinnerAdapter) Warning(text string) {
	p.spinner.Warning(text)
}

func (p *ptermSpinnerAdapter) Stop() error {
	return p.spinner.Stop()
}

// StartSpinner creates and starts a spinner if appropriate for the environment
func (om *OutputManager) StartSpinner(spinner pterm.SpinnerPrinter, text string, writer ...*pterm.MultiPrinter) SpinnerHandle {
	if !om.config.ShowProgress {
		return &NoOpSpinner{}
	}

	var ps *pterm.SpinnerPrinter
	var err error

	if len(writer) > 0 && writer[0] != nil {
		ps, err = spinner.WithWriter(writer[0].NewWriter()).Start(text)
	} else {
		ps, err = spinner.Start(text)
	}

	if err != nil {
		return &NoOpSpinner{}
	}

	return &ptermSpinnerAdapter{spinner: ps}
}

// ShouldShowProgress returns true if progress indicators should be shown
func (om *OutputManager) ShouldShowProgress() bool {
	return om.config.ShowProgress
}

// IsVerbose returns true if verbose output is enabled
func (om *OutputManager) IsVerbose() bool {
	return om.config.Verbose
}

// IsSilent returns true if silent mode is enabled
func (om *OutputManager) IsSilent() bool {
	return om.config.Silent
}

// ShouldGenerateReport returns true if report should be generated
func (om *OutputManager) ShouldGenerateReport() bool {
	return om.config.GenerateReport
}

// Print outputs a message if not in silent mode
func (om *OutputManager) Print(a ...interface{}) {
	if !om.config.Silent {
		pterm.Print(a...)
	}
}

// Println outputs a message with newline if not in silent mode
func (om *OutputManager) Println(a ...interface{}) {
	if !om.config.Silent {
		pterm.Println(a...)
	}
}

// Printf outputs a formatted message if not in silent mode
func (om *OutputManager) Printf(format string, a ...interface{}) {
	if !om.config.Silent {
		pterm.Printf(format, a...)
	}
}

// VerbosePrint outputs a message only in verbose mode
func (om *OutputManager) VerbosePrint(a ...interface{}) {
	if om.config.Verbose {
		pterm.Print(a...)
	}
}

// VerbosePrintln outputs a message with newline only in verbose mode
func (om *OutputManager) VerbosePrintln(a ...interface{}) {
	if om.config.Verbose {
		pterm.Println(a...)
	}
}

// VerbosePrintf outputs a formatted message only in verbose mode
func (om *OutputManager) VerbosePrintf(format string, a ...interface{}) {
	if om.config.Verbose {
		pterm.Printf(format, a...)
	}
}
