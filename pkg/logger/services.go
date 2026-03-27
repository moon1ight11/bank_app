package logger

func (l *appLogger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

func (l *appLogger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *appLogger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

func (l *appLogger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}
