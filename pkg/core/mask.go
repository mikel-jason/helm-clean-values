package core

import (
	"reflect"

	"golang.org/x/exp/slog"
)

type valueStack struct {
	identifier     string
	fullIdentifier []string
	value          interface{}
}

type compareResult struct {
	stack   valueStack
	finding int
	childs  []compareResult
	keep    bool
}

const (
	compareResultTypeMismatch = iota
	compareResultNotImplemented
	compareResultTypeMatch
	compareResultDoesNotExistOnReference
)

func Mask(input, reference map[string]interface{}) map[string]interface{} {
	inputStack := valueStack{
		"", []string{}, input,
	}
	referenceStack := valueStack{
		"", []string{}, reference,
	}

	res := compare(inputStack, referenceStack)
	out := write(res)
	if outMap, ok := out.(map[string]interface{}); ok { // cannot cast to map if empty
		return outMap
	}
	return make(map[string]interface{})
}

func compare(input, reference valueStack) compareResult {

	inputType := reflect.ValueOf(input.value).Kind()
	referenceType := reflect.ValueOf(reference.value).Kind()

	if inputType != referenceType {
		if referenceType == reflect.Invalid {
			slog.Info("found non-existing", "id", input.fullIdentifier)
			return compareResult{input, compareResultDoesNotExistOnReference, []compareResult{}, false}
		}
		slog.Info("found type mismatch", "id", input.fullIdentifier)
		return compareResult{input, compareResultTypeMismatch, []compareResult{}, false}
	}

	switch inputType {
	case reflect.Map:
		slog.Info("found map", "id", input.fullIdentifier)
		childs := []compareResult{}
		keep := false
		for key, val := range input.value.(map[string]interface{}) {
			res := compare(valueStack{
				identifier:     key,
				fullIdentifier: append(input.fullIdentifier, key),
				value:          val,
			}, valueStack{
				identifier:     key,
				fullIdentifier: append(reference.fullIdentifier, key),
				value:          reference.value.(map[string]interface{})[key], // must be map, type mismatch is filtered before
			})
			if res.keep {
				keep = true
			}
			childs = append(childs, res)
		}
		return compareResult{input, compareResultTypeMatch, childs, keep}
	case reflect.String, reflect.Bool, reflect.Int:
		slog.Info("found string/bool/int", "id", input.fullIdentifier)
		return compareResult{input, compareResultTypeMatch, []compareResult{}, true}
	case reflect.Slice:
		slog.Info("found slice", "id", input.fullIdentifier)
		return compareResult{input, compareResultNotImplemented, []compareResult{}, true}
	default:
		slog.Info("found unknown", "id", input.fullIdentifier)
		return compareResult{input, compareResultNotImplemented, []compareResult{}, false}
	}
}

func write(input compareResult) interface{} {

	if !input.keep {
		return nil
	}

	if len(input.childs) == 0 {
		return input.stack.value
	}

	res := make(map[string]interface{})
	for _, child := range input.childs {
		val := write(child)
		if val != nil {
			res[child.stack.identifier] = val
		}
	}
	return res
}
