package file

import (
	"fmt"
	"os"

	"github.com/sarcaustech/helm-clean-values/pkg/core"
	"gopkg.in/yaml.v3"
)

type ValuesProvider struct {
	Path string
}

func (p *ValuesProvider) Values(core.Logger) (map[string]interface{}, error) {
	bytes, err := os.ReadFile(p.Path)
	if err != nil {
		return nil, fmt.Errorf("cannot not read file %s: %w", p.Path, err)
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("cannot parse file contents to YAML %s: %w", p.Path, err)
	}
	return data, nil
}
