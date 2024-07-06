package record

func (r *Record) FileName() string {
	return r.file
}

func (r *Record) LogType() string {
	return r.logType
}

func (r *Record) Message() string {
	return r.message
}
