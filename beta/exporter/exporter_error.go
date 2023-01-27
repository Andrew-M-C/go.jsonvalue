package exporter

// E is the internal error type of package exporter.
type E string

// Error implements error interface
func (e E) Error() string {
	return string(e)
}

const (
	ErrUnsupportedType = E("jsonvalue.exporter: unsupported type")
)
