package module

type Type string

const (
	DOWNLOADER Type = "downloader"
	ANALYZER   Type = "analyzer"
	PIPELINE   Type = "pipeline"
)

var LegalTypeLetterMap = map[Type]string{
	DOWNLOADER: "D",
	ANALYZER:   "A",
	PIPELINE:   "P",
}

var LegalLetterTypeMap = map[string]Type{
	"D": DOWNLOADER,
	"A": ANALYZER,
	"P": PIPELINE,
}
