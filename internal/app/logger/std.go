package logger

var (
	std Logger = NewLogrus()
)

func SetStdLogger(l Logger) {
	std = l
}

func StdLogger() Logger {
	return std
}
