package stdin

import (
	"fmt"
	"io"
	"os"

	"github.com/sarcaustech/helm-clean-values/pkg/core"
	"gopkg.in/yaml.v3"
)

type ValuesProvider struct{}

func (p *ValuesProvider) Values(core.Logger) (map[string]interface{}, error) {
	valuesBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("cannot read stdin: %w", err)
	}

	var values map[string]interface{}
	if err = yaml.Unmarshal(valuesBytes, &values); err != nil {
		return nil, fmt.Errorf("cannot parse stdin to YAML: %w", err)
	}

	return values, nil
}
