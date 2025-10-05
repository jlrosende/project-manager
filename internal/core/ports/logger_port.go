package ports

// LogField represents a structured logging attribute passed across layers.
type LogField struct {
	Key   string
	Value any
}

// Logger abstracts structured logging so application services avoid depending
// on concrete logging frameworks.
type Logger interface {
	Debug(msg string, fields ...LogField)
	Info(msg string, fields ...LogField)
	Warn(msg string, fields ...LogField)
	Error(msg string, fields ...LogField)
}
