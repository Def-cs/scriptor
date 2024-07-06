package event

func (l *Event) Name() string {
	return l.name
}

func (l *Event) Prefix() string {
	return l.prefix
}

func (l *Event) Flag() int {
	return l.flag
}
