package log

type NoopLogger struct {
}

func (n NoopLogger) Debug(msg string, args ...any) {
}
