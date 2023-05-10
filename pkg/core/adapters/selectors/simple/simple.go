package simple

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sarcaustech/helm-clean-values/pkg/core"
	"golang.org/x/exp/slog"
)

type Selector struct{}

func (s *Selector) Run(logger core.Logger, input, reference map[string]interface{}) (core.SelectResult, error) {
	result := compare(logger, simpleResult{
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

func compare(logger core.Logger, elem simpleResult) simpleResult {

	inputType := reflect.ValueOf(elem.inputValue).Kind()
	referenceType := reflect.ValueOf(elem.referenceValue).Kind()

	if inputType != referenceType {
		if referenceType == reflect.Invalid {
			logger.Debug(fmt.Sprintf("found non-existing, id: %s", strings.Join(elem.fullIdentifier, ".")))
			elem.keep = false
			elem.reason = core.ReasonDoesNotExistOnReference
			return elem
		}
		logger.Debug(fmt.Sprintf("found type missmatch, id: %s - user input is %s, but chart has %s",
			strings.Join(elem.fullIdentifier, "."), inputType.String(), referenceType.String()))
		elem.keep = false
		elem.reason = core.ReasonTypeMismatch
		return elem
	}

	switch inputType {
	case reflect.Map:
		logger.Debug(fmt.Sprintf("found map, id: %s", strings.Join(elem.fullIdentifier, ".")))
		elem.keep = false
		asMap, ok := elem.inputValue.(map[string]interface{})
		if !ok { // nok w/ gopkg.in/yaml.v2
			slog.Warn("found map with unknown key type", "id", elem.fullIdentifier, "value", elem.inputValue)
			logger.Warn(fmt.Sprintf("found map with unknown key type, id: %s, kind: %s, type: %s",
				strings.Join(elem.fullIdentifier, "."), inputType.String(), reflect.ValueOf(elem.inputValue).Type().String()))
			elem.keep = true
			elem.reason = core.ReasonNotImplemented

			return elem
		}
		for key, val := range asMap {

			res := compare(logger, simpleResult{
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
		logger.Debug(fmt.Sprintf("found known primitive string/bool/int, id: %s", strings.Join(elem.fullIdentifier, ".")))
		elem.keep = true
		elem.reason = core.ReasonTypeMatch
		return elem
	case reflect.Slice:
		logger.Debug(fmt.Sprintf("found list, id: %s", strings.Join(elem.fullIdentifier, ".")))
		refAsList := elem.referenceValue.([]interface{})
		mergedRef, err := mergeSlice(refAsList)
		if err != nil {
			logger.Error(fmt.Sprintf("error while merging list %s: %s", strings.Join(elem.fullIdentifier, "."), err.Error()))
		}

		for idx, listItem := range elem.inputValue.([]interface{}) {
			logger.Debug(fmt.Sprintf("comparing list element: %v | %v", listItem, refAsList))
			key := strconv.Itoa(idx)
			res := compare(logger, simpleResult{
				localIdentifier: key,
				fullIdentifier:  append(elem.fullIdentifier, key),
				inputValue:      listItem,
				referenceValue:  mergedRef,
			})

			if res.keep { // if any child is ok, keep parent
				elem.keep = true
			}

			elem.childs = append(elem.childs, res)
		}

		return elem
	default:
		logger.Debug(fmt.Sprintf("found unknown type, will drop, id: %s", strings.Join(elem.fullIdentifier, ".")))
		elem.keep = false
		elem.reason = core.ReasonNotImplemented
		return elem
	}
}

func mergeSlice(input []interface{}) (interface{}, error) {
	if len(input) == 0 {
		return []interface{}{}, nil
	}
	zeroKind := reflect.ValueOf(input[0]).Kind()

	if zeroKind != reflect.Slice { // test outside of loop for performance
		var res interface{}
		for _, child := range input {
			if zeroKind != reflect.ValueOf(child).Kind() {
				return nil, fmt.Errorf("cannot merge list of different types, mismating elements: %v | %v", input[0], child)
			}
			subRes, err := merge(res, child)
			if err != nil {
				return nil, err // exit early for performance
			}
			res = subRes
		}
		return res, nil
	}

	return nil, fmt.Errorf("list of lists not implemented yet")
}

// merge merges two variables based on their types. It does not care about values,
// only builds e.g. a full map[string]interface{} that structurally has all child
// elements of both variables. It prefers the left values if both exist.
func merge(left, right interface{}) (interface{}, error) {
	if left == nil {
		return right, nil // can be nil
	} else if right == nil {
		return left, nil
	}

	leftKind := reflect.ValueOf(left).Kind()
	rightKind := reflect.ValueOf(left).Kind()

	if leftKind != rightKind {
		return nil, fmt.Errorf("values do not have the same type: %v | %v", left, right)
	}

	switch leftKind {
	case reflect.Map:
		res := map[string]interface{}{}
		leftMap := left.(map[string]interface{})
		rightMap := right.(map[string]interface{})
		for keyLeft, valueLeft := range leftMap {
			valueRight, existsRight := rightMap[keyLeft]
			if !existsRight {
				res[keyLeft] = valueLeft
			} else {
				childRes, err := merge(valueLeft, valueRight)
				if err != nil {
					return nil, err // exit early for performance
				}
				res[keyLeft] = childRes
			}
		}

		for rightKey, rightValue := range rightMap {
			if _, alreadyMerged := res[rightKey]; alreadyMerged {
				continue
			} // if it does not exist yet, it cannot exist in leftMap
			res[rightKey] = rightValue
		}

		return res, nil

	case reflect.Slice:
		leftSlice := left.([]interface{})
		rightSlice := right.([]interface{})

		leftMerge := []interface{}{}
		for _, elem := range leftSlice {
			subRes, err := merge(leftMerge, elem)
			if err != nil {
				return nil, err // exit early for performance
			}
			leftMerge = subRes.([]interface{})
		}

		rightMerge := []interface{}{}
		for _, elem := range rightSlice {
			subRes, err := merge(rightMerge, elem)
			if err != nil {
				return nil, err // exit early for performance
			}
			rightMerge = subRes.([]interface{})
		}

		return merge(leftMerge, rightMerge)
	default:
		return left, nil
	}
}
