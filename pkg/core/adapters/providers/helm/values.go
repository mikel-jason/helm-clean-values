package helm

import (
	"fmt"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type ValuesProvider struct {
	Prompt     string
	BinaryPath string
}

func (p *ValuesProvider) Values() (map[string]interface{}, error) {
	cmd := exec.Command(p.BinaryPath, "show", "values", p.Prompt)
	valuesBytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("helm show values %s failed: %w", p.Prompt, err)
	}

	var values map[string]interface{}
	if err = yaml.Unmarshal(valuesBytes, &values); err != nil {
		return nil, fmt.Errorf("cannot parse stdin to YAML: %w", err)
	}

	return values, nil
}
