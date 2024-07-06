package record

type Record struct {
	file    string
	message string
	logType string
}

func NewRecord(file string, message string, t string) *Record {
	return &Record{
		file:    file,
		message: message,
		logType: t,
	}
}
