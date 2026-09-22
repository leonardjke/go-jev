package log

type NoopLogger struct{}

func (NoopLogger) Debug(_ string, _ ...any) {}
