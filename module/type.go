package module

type Type string

const (
	DOWNLOADER Type = "downloader"
	ANALYZER   Type = "analyzer"
	PIPELINE   Type = "pipeline"
)

var legalTypeLetterMap = map[Type]string{
	DOWNLOADER: "D",
	ANALYZER:   "A",
	PIPELINE:   "P",
}
