package core

type ValuesProvider interface {
	Values() (map[string]interface{}, error)
}

type ValueSelector interface {
	Run(input, reference map[string]interface{}) (SelectResult, error)
}

type SelectResult struct {
	LocalIdentifier string
	FullIdentifier  []string
	Value           interface{}
	Keep            bool
	Reason          int
	Childs          []SelectResult
}

const (
	ReasonUndefined = iota // default int value 0 -> implicitly set
	ReasonTypeMatch
	ReasonTypeMismatch
	ReasonDoesNotExistOnReference
	ReasonEmpty
	ReasonNotImplemented
)
