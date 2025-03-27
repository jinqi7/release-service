package resty

type NullLogger struct{}

func (n *NullLogger) Errorf(format string, v ...interface{}) {
}

func (n *NullLogger) Warnf(format string, v ...interface{}) {
}

func (n *NullLogger) Debugf(format string, v ...interface{}) {
}
