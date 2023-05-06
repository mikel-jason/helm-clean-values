package simple

import (
	"reflect"

	"github.com/sarcaustech/helm-clean-values/pkg/core"
	"golang.org/x/exp/slog"
)

type Selector struct{}

func (s *Selector) Run(input, reference map[string]interface{}) (core.SelectResult, error) {
	result := compare(simpleResult{
		localIdentifier: "",
		fullIdentifier:  []string{},
		inputValue:      input,
		referenceValue:  reference,
	})

	return result.coreResult(), nil

}

type simpleResult struct {
	localIdentifier string
	fullIdentifier  []string
	inputValue      interface{}
	referenceValue  interface{}
	keep            bool
	reason          int
	childs          []simpleResult
}

func (s *simpleResult) coreResult() core.SelectResult {
	childs := []core.SelectResult{}
	for _, c := range s.childs {
		childs = append(childs, c.coreResult())
	}
	return core.SelectResult{
		LocalIdentifier: s.localIdentifier,
		FullIdentifier:  s.fullIdentifier,
		Value:           s.inputValue,
		Keep:            s.keep,
		Reason:          s.reason,
		Childs:          childs,
	}
}

func compare(elem simpleResult) simpleResult {

	inputType := reflect.ValueOf(elem.inputValue).Kind()
	referenceType := reflect.ValueOf(elem.referenceValue).Kind()

	if inputType != referenceType {
		if referenceType == reflect.Invalid {
			slog.Info("found non-existing", "id", elem.fullIdentifier)
			elem.keep = false
			elem.reason = core.ReasonDoesNotExistOnReference
			return elem
		}
		slog.Info("found type mismatch", "id", elem.fullIdentifier)
		elem.keep = false
		elem.reason = core.ReasonTypeMismatch
		return elem
	}

	switch inputType {
	case reflect.Map:
		slog.Info("found map", "id", elem.fullIdentifier)
		elem.keep = false
		asMap, ok := elem.inputValue.(map[string]interface{})
		if !ok { // nok w/ gopkg.in/yaml.v2
			slog.Warn("found map with unknown key type", "id", elem.fullIdentifier, "value", elem.inputValue)
			elem.keep = true
			elem.reason = core.ReasonNotImplemented

			return elem
		}
		for key, val := range asMap {

			res := compare(simpleResult{
				localIdentifier: key,
				fullIdentifier:  append(elem.fullIdentifier, key),
				inputValue:      val,
				referenceValue:  elem.referenceValue.(map[string]interface{})[key], // must be map, type mismatch is filtered before
			})

			if res.keep { // if any child is ok, keep parent
				elem.keep = true
			}

			elem.childs = append(elem.childs, res)
		}
		if elem.keep {
			elem.reason = core.ReasonTypeMatch
		} else {
			elem.reason = core.ReasonEmpty
		}
		return elem
	case reflect.String, reflect.Bool, reflect.Int:
		slog.Info("found string/bool/int", "id", elem.fullIdentifier)
		elem.keep = true
		elem.reason = core.ReasonTypeMatch
		return elem
	case reflect.Slice:
		slog.Info("found slice", "id", elem.fullIdentifier)
		elem.keep = true
		elem.reason = core.ReasonNotImplemented
		return elem
	default:
		slog.Info("found unknown", "id", elem.fullIdentifier)
		elem.keep = false
		elem.reason = core.ReasonNotImplemented
		return elem
	}
}
