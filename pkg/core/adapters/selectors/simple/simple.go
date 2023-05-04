package simple

import (
	"reflect"

	"github.com/sarcaustech/helm-clean-values/pkg/core"
	"golang.org/x/exp/slog"
)

type Selector struct{}

func (s *Selector) Run(input, reference map[string]interface{}) (core.SelectResult, error) {
	return compare(core.SelectResult{
		LocalIdentifier: "",
		FullIdentifier:  []string{},
		InputValue:      input,
		ReferenceValue:  reference,
	}), nil
}

func compare(elem core.SelectResult) core.SelectResult {

	inputType := reflect.ValueOf(elem.InputValue).Kind()
	referenceType := reflect.ValueOf(elem.ReferenceValue).Kind()

	if inputType != referenceType {
		if referenceType == reflect.Invalid {
			slog.Info("found non-existing", "id", elem.FullIdentifier)
			elem.Keep = false
			elem.Reason = core.ReasonDoesNotExistOnReference
			return elem
		}
		slog.Info("found type mismatch", "id", elem.FullIdentifier)
		elem.Keep = false
		elem.Reason = core.ReasonTypeMismatch
		return elem
	}

	switch inputType {
	case reflect.Map:
		slog.Info("found map", "id", elem.FullIdentifier)
		elem.Keep = false
		asMap, ok := elem.InputValue.(map[string]interface{})
		if !ok { // nok w/ gopkg.in/yaml.v2
			slog.Warn("found map with unknown key type", "id", elem.FullIdentifier, "value", elem.InputValue)
			elem.Keep = true
			elem.Reason = core.ReasonNotImplemented

			return elem
		}
		for key, val := range asMap {

			res := compare(core.SelectResult{
				LocalIdentifier: key,
				FullIdentifier:  append(elem.FullIdentifier, key),
				InputValue:      val,
				ReferenceValue:  elem.ReferenceValue.(map[string]interface{})[key], // must be map, type mismatch is filtered before
			})

			if res.Keep { // if any child is ok, keep parent
				elem.Keep = true
			}

			elem.Childs = append(elem.Childs, res)
		}
		if elem.Keep {
			elem.Reason = core.ReasonTypeMatch
		} else {
			elem.Reason = core.ReasonEmpty
		}
		return elem
	case reflect.String, reflect.Bool, reflect.Int:
		slog.Info("found string/bool/int", "id", elem.FullIdentifier)
		elem.Keep = true
		elem.Reason = core.ReasonTypeMatch
		return elem
	case reflect.Slice:
		slog.Info("found slice", "id", elem.FullIdentifier)
		elem.Keep = true
		elem.Reason = core.ReasonNotImplemented
		return elem
	default:
		slog.Info("found unknown", "id", elem.FullIdentifier)
		elem.Keep = false
		elem.Reason = core.ReasonNotImplemented
		return elem
	}
}
