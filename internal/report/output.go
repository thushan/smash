package report

type ReportOutput struct {
	Summary    *RunSummary
	Duplicates []SmashFile
}
