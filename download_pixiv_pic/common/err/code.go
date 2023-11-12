package err

const (
	Success      = 200
	Error        = 500
	InvalidParam = 400

	RootError          int = 601
	ConfigFileNotFound     = 00001
	ConfigFileReadErr      = 00002
	ConfigReadErr          = 00003
	ConfigReadSuccess      = 00004
)
