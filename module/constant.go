package module

type ErrorType string

const (
	DownloaderError ErrorType = "downloader error"
	AnalyzerError   ErrorType = "analyzer error"
	PipelineError   ErrorType = "pipeline error"
	SchedulerError  ErrorType = "scheduler error"
)
