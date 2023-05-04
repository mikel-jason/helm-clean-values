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
	InputValue      interface{}
	ReferenceValue  interface{}
	Keep            bool
	Reason          int
	Childs          []SelectResult
}

const (
	ReasonTypeMatch = iota
	ReasonTypeMismatch
	ReasonDoesNotExistOnReference
	ReasonEmpty
	ReasonNotImplemented
)
