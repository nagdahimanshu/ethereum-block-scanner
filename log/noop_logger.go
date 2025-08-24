package log

type noOpLogger struct{}

func (l *noOpLogger) Debug(...interface{})                     {}
func (l *noOpLogger) Debugf(msg string, others ...interface{}) {}
func (l *noOpLogger) Debugw(msg string, others ...interface{}) {}
func (l *noOpLogger) Info(...interface{})                      {}
func (l *noOpLogger) Infof(msg string, others ...interface{})  {}
func (l *noOpLogger) Infow(msg string, others ...interface{})  {}
func (l *noOpLogger) Warn(...interface{})                      {}
func (l *noOpLogger) Warnf(msg string, others ...interface{})  {}
func (l *noOpLogger) Warnw(msg string, others ...interface{})  {}
func (l *noOpLogger) Error(...interface{})                     {}
func (l *noOpLogger) Errorf(msg string, others ...interface{}) {}
func (l *noOpLogger) Errorw(msg string, others ...interface{}) {}
func (l *noOpLogger) Fatal(...interface{})                     {}
func (l *noOpLogger) Fatalf(msg string, others ...interface{}) {}
func (l *noOpLogger) Fatalw(msg string, others ...interface{}) {}
func (l *noOpLogger) With(kv ...interface{}) Logger {
	return &noOpLogger{}
}

func NewNoOpLogger() Logger {
	return &noOpLogger{}
}

// Ensure noOpLogger conforms to the Logger interface.
var _ Logger = (*noOpLogger)(nil)
