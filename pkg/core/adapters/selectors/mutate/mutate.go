package mutate

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/sarcaustech/helm-clean-values/pkg/core"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

type Selector struct {
	HelmBinaryPath string
	Prompt         string
	originalResult []byte // template result of input, to be tested against
	bytesDeepCopy  []byte // the input yaml as bytes for re-parsing as deep copy
}

func (s *Selector) Run(logger core.Logger, input, reference map[string]interface{}) (core.SelectResult, error) {

	var err error
	s.originalResult, err = s.template(input)
	if err != nil {
		return core.SelectResult{}, fmt.Errorf("error templating with original user data: %w", err)
	}

	s.bytesDeepCopy, err = yaml.Marshal(input)
	if err != nil {
		panic(err) // string -> typed -> string should always succeed
	}

	prep := s.prepare(logger, mutateResult{
		Remaining: input,
	})

	return prep.decide(logger, s), nil
}

func (s *Selector) template(values map[string]interface{}) ([]byte, error) {

	yamlValuesBytes, err := yaml.Marshal(values)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot parse values to YAML: %w", err)
	}

	cmd := exec.Command(s.HelmBinaryPath, "template", s.Prompt, "-f", "-")
	cmdIn, err := cmd.StdinPipe()
	if err != nil {
		return []byte{}, fmt.Errorf("cannot open helm template STDIN: %w", err)
	}

	_, err = io.WriteString(cmdIn, string(yamlValuesBytes))
	if err != nil {
		cmdIn.Close()
		return []byte{}, fmt.Errorf("cannot write to helm template STDIN: %w", err)
	}

	cmdIn.Close()

	outBytes, err := cmd.Output()
	if err != nil {
		return []byte{}, fmt.Errorf("cannot template the chart: %w", err)
	}

	return outBytes, nil
}

type mutateResult struct {
	Local     string
	Path      []string
	Remaining interface{}
	Childs    []mutateResult
	Keep      bool
}

func (m *mutateResult) decide(logger core.Logger, s *Selector) core.SelectResult {
	res := core.SelectResult{
		LocalIdentifier: m.Local,
		FullIdentifier:  m.Path,
		Value:           m.Remaining,
		Reason:          core.ReasonUndefined,
	}

	for _, c := range m.Childs {
		childRes := c.decide(logger, s)
		res.Childs = append(res.Childs, childRes)

		if childRes.Keep {
			res.Keep = true
		}
	}

	if !res.Keep && len(res.Childs) == 0 {
		var mutatedMap map[string]interface{}
		err := yaml.Unmarshal(s.bytesDeepCopy, &mutatedMap)
		if err != nil {
			panic(err)
		}

		err = setValueByPath(logger, mutatedMap, m.Path, mutatePrimitive)
		if err != nil {
			slog.Error("cannot mutate primitive in user values, skipping and assuming as relevant",
				"path", strings.Join(m.Path, "."))
		} else {
			result, err := s.template(mutatedMap)
			if err != nil {
				panic(err)
			}

			res.Keep = !bytes.Equal(result, s.originalResult)
		}

	}

	/*
		TODO
		built-in assumption:
		a map is not needed if none of its childs are needed, but charts could loop maps and ignore
		the values, so keys only can be relevant, too.

		implement this edge case
	*/

	return res
}

func (s *Selector) prepare(logger core.Logger, input mutateResult) mutateResult {
	remainingKind := reflect.ValueOf(input.Remaining).Kind()

	switch remainingKind {
	case reflect.Map:
		logger.Debug(fmt.Sprintf("found map, id: %s", strings.Join(input.Path, ".")))
		for key, value := range input.Remaining.(map[string]interface{}) {
			childResult := s.prepare(logger, mutateResult{
				Local:     key,
				Path:      append(input.Path, key),
				Remaining: value,
			})

			input.Childs = append(input.Childs, childResult)
		}
	case reflect.Slice:
		logger.Warn(fmt.Sprintf("found list which is not supported yet, id: %s", strings.Join(input.Path, ".")))
		return input
	default:
		logger.Debug(fmt.Sprintf("found primitive, id: %s", strings.Join(input.Path, ".")))
		return input
	}
	return input
}

func mutatePrimitive(logger core.Logger, input interface{}) interface{} {
	switch reflect.ValueOf(input).Kind() {
	case reflect.String:
		return input.(string) + "a"
	case reflect.Bool:
		return !input.(bool)
	case reflect.Int:
		return input.(int) + 1
	default:
		logger.Info(fmt.Sprintf("unknown kind, replacing with nil, value: %v", input))
		return nil
	}
}

func setValueByPath(logger core.Logger, data interface{}, path []string, op func(core.Logger, interface{}) interface{}) error {
	if len(path) == 0 {
		return nil
	}
	for _, field := range path[:len(path)-1] {
		if typedData, ok := data.(map[string]interface{}); ok {
			data, ok = typedData[field]
			if !ok {
				return fmt.Errorf("path not found: %s", field)
			}
		} else if typedData, ok := data.([]interface{}); ok {
			idx, err := strconv.Atoi(field)
			if err != nil {
				return fmt.Errorf("path not found: %s", field)
			}
			if !ok || idx >= len(typedData) {
				return fmt.Errorf("path not found: %s", field)
			}
			data = typedData[idx]
		} else {
			return fmt.Errorf("invalid YAML data type")
		}
	}
	lastField := path[len(path)-1]

	if typedData, ok := data.(map[string]interface{}); ok {
		typedData[lastField] = op(logger, typedData[lastField])
	} else if typedData, ok := data.([]interface{}); ok {
		idx, err := strconv.Atoi(lastField)
		if err != nil {
			return fmt.Errorf("path not found: %s", lastField)
		}
		if !ok || idx >= len(typedData) {
			return fmt.Errorf("path not found: %s", lastField)
		}
		typedData[idx] = op(logger, typedData[idx])
	} else {
		return fmt.Errorf("invalid YAML data type")
	}
	return nil
}
