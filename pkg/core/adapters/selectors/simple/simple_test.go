package simple_test

import (
	"testing"

	"github.com/sarcaustech/helm-clean-values/pkg/core"
	"github.com/sarcaustech/helm-clean-values/pkg/core/adapters/selectors/simple"
	"github.com/stretchr/testify/require"
)

func TestValues(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		reference map[string]interface{}
		expected  map[string]interface{}
	}{
		{
			name: "flat fit",
			input: map[string]interface{}{
				"string": "hello",
				"int":    123,
				"bool":   true,
			},
			reference: map[string]interface{}{
				"string":     "world",
				"int":        321,
				"bool":       false,
				"additional": "unused",
			},
			expected: map[string]interface{}{
				"string": "hello",
				"int":    123,
				"bool":   true,
			},
		},
		{
			name: "nested fit",
			input: map[string]interface{}{
				"nested": map[string]interface{}{
					"map": 1,
				},
			},
			reference: map[string]interface{}{
				"string": "unused",
				"nested": map[string]interface{}{
					"additional": "unused",
					"map":        1,
				},
			},
			expected: map[string]interface{}{
				"nested": map[string]interface{}{
					"map": 1,
				},
			},
		},
		{
			name:  "empty input",
			input: map[string]interface{}{},
			reference: map[string]interface{}{
				"string": "unused",
			},
			expected: map[string]interface{}{},
		},
		{
			name: "empty reference",
			input: map[string]interface{}{
				"string": "unused",
			},
			reference: map[string]interface{}{},
			expected:  map[string]interface{}{},
		},
		{
			name: "flat drop",
			input: map[string]interface{}{
				"string": "hello",
				"int":    123,
				"bool":   true,
			},
			reference: map[string]interface{}{
				"string": "world",
			},
			expected: map[string]interface{}{
				"string": "hello",
			},
		},
		{
			name: "nested drop",
			input: map[string]interface{}{
				"nested": map[string]interface{}{
					"map":        1,
					"additional": "unused",
				},
				"additional": "unused",
			},
			reference: map[string]interface{}{
				"string": "unused",
				"nested": map[string]interface{}{
					"map": 1,
				},
			},
			expected: map[string]interface{}{
				"nested": map[string]interface{}{
					"map": 1,
				},
			},
		},
	}

	require := require.New(t)
	selector := simple.Selector{}
	for _, test := range tests {
		result, err := selector.Run(test.input, test.reference)
		require.Nilf(err, test.name)
		clean, err := core.Populate(result)
		require.Nilf(err, test.name)
		require.Equalf(test.expected, clean, test.name)
	}
}
