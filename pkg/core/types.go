package core

type ValuesProvider interface {
	Values(Logger) (map[string]interface{}, error)
}

type ValueSelector interface {
	Run(Logger, map[string]interface{}, map[string]interface{}) (SelectResult, error)
}

type SelectResult struct {
	LocalIdentifier string
	FullIdentifier  []string
	Value           interface{}
	Keep            bool
	Reason          int
	Childs          []SelectResult
}

type Logger interface {
	Error(string)
	Warn(string)
	Info(string)
	Debug(string)
}

const (
	ReasonUndefined = iota // default int value 0 -> implicitly set
	ReasonTypeMatch
	ReasonTypeMismatch
	ReasonDoesNotExistOnReference
	ReasonEmpty
	ReasonNotImplemented
)
