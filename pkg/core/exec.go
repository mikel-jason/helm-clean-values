package core

import (
	"fmt"
)

func Run(
	inputProvider ValuesProvider,
	referenceProvider ValuesProvider,
	selector ValueSelector) (map[string]interface{}, error) {
	input, err := inputProvider.Values()
	if err != nil {
		return nil, fmt.Errorf("cannot get input values: %w", err)
	}
	reference, err := referenceProvider.Values()
	if err != nil {
		return nil, fmt.Errorf("cannot get reference values: %w", err)
	}

	selects, err := selector.Run(input, reference)
	if err != nil {
		return nil, fmt.Errorf("error during value selections: %w", err)
	}

	checkKeepFromChilds(&selects)

	cleanValues, err := Populate(selects)
	if err != nil {
		return nil, err
	}
	return cleanValues, nil
}

func Populate(input SelectResult) (map[string]interface{}, error) {
	res, ok := populate(input).(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid dict %v", res)
	}
	return res, nil
}

func populate(input SelectResult) interface{} {

	if !input.Keep && len(input.FullIdentifier) == 0 {
		return map[string]interface{}{}
	}

	if !input.Keep {
		return nil
	}

	if len(input.Childs) == 0 {
		return input.Value
	}

	res := make(map[string]interface{})
	for _, child := range input.Childs {
		val := populate(child)
		if val != nil {
			res[child.LocalIdentifier] = val
		}
	}
	return res
}

// checkKeepFromChilds traverses a SelectResult tree and ensures that no node it not kept that has a (nested) child to be kept
func checkKeepFromChilds(input *SelectResult) {
	if input.Keep {
		return
	}

	if len(input.Childs) == 0 {
		return
	}

	for _, child := range input.Childs {
		checkKeepFromChilds(&child)
		if child.Keep {
			input.Keep = true
			return
		}
	}
}
