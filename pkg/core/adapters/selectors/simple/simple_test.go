package simple

import (
	"encoding/json"
	"testing"

	"github.com/sarcaustech/helm-clean-values/pkg/core"
	"github.com/sarcaustech/helm-clean-values/pkg/logger"
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
	selector := Selector{}
	logger := &logger.Plain{EnableDebug: true}
	for _, test := range tests {
		result, err := selector.Run(logger, test.input, test.reference)
		require.Nilf(err, test.name)
		clean, err := core.Populate(result)
		require.Nilf(err, test.name)
		require.Equalf(test.expected, clean, test.name)
	}
}

func TestMergeSlice(t *testing.T) {
	require := require.New(t)

	tests := []struct {
		title string
		input []interface{}
		// input       func()([]interface{})
		expected    interface{}
		expectedErr bool
	}{
		{
			title:    "flat integers",
			input:    []interface{}{1, 2, 3},
			expected: 1,
		},
		{
			title:    "flat strings",
			input:    []interface{}{"1", "2", "3"},
			expected: "1",
		},
		{
			title: "maps",
			input: []interface{}{
				map[string]interface{}{
					"first":    1,
					"firstTwo": 1,
					"all":      "first",
				},
				map[string]interface{}{
					"second":   2,
					"firstTwo": 2,
					"lastTwo":  2,
					"all":      "second",
				},
				map[string]interface{}{
					"third":   3,
					"lastTwo": 3,
					"all":     "second",
				},
			},
			expected: map[string]interface{}{
				"first":    1,
				"second":   2,
				"third":    3,
				"firstTwo": 1,
				"lastTwo":  2,
				"all":      "first",
			},
		},
		{
			title: "list of maps",
			input: []interface{}{
				[]map[string]interface{}{
					{
						"first": 1,
						"list":  1,
						"all":   "first@1",
					},
					{
						"second": 2,
						"list":   2,
						"all":    "second@1",
					},
				},
				[]map[string]interface{}{
					{
						"third": 3,
						"list2": 3,
						"all":   "third@2",
					},
				},
			},
			expectedErr: true,
			// expected: ???,
		},
	}

	for _, test := range tests {
		expectedJson, err := json.Marshal(test.expected)
		if !test.expectedErr && err != nil {
			require.Nilf(err, "test case setup wrong, cannot parse expected to json", test.title)
		}

		result, err := mergeSlice(test.input)
		if test.expectedErr {
			require.NotNilf(err, test.title)
			continue
		}
		require.Nilf(err, test.title)
		resultJson, err := json.Marshal(result)
		require.Nilf(err, test.title)

		require.Equalf(expectedJson, resultJson, test.title)

	}
}

func TestMerge(t *testing.T) {
	require := require.New(t)

	tests := []struct {
		title       string
		left        interface{}
		right       interface{}
		expected    interface{}
		expectedErr bool
	}{
		{
			title:    "two strings",
			left:     "a",
			right:    "b",
			expected: "a",
		},
		{
			title:    "two ints",
			left:     1,
			right:    2,
			expected: 1,
		},
		{
			title:    "two bools",
			left:     true,
			right:    false,
			expected: true,
		},
		{
			title: "two flat maps",
			left: map[string]interface{}{
				"left": 1,
				"both": "left",
			},
			right: map[string]interface{}{
				"right": 1,
				"both":  "right",
			},
			expected: map[string]interface{}{
				"right": 1,
				"left":  1,
				"both":  "left",
			},
		},
	}

	for _, test := range tests {
		expectedJson, err := json.Marshal(test.expected)
		if !test.expectedErr && err != nil {
			require.Nilf(err, "test case setup wrong, cannot parse expected to json", test.title)
		}

		result, err := merge(test.left, test.right)
		if test.expectedErr {
			require.NotNilf(err, test.title)
			continue
		}
		require.Nilf(err, test.title)
		resultJson, err := json.Marshal(result)
		require.Nilf(err, test.title)

		require.Equalf(expectedJson, resultJson, test.title)

	}
}
