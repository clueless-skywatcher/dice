package commands

import (
	"testing"

	"github.com/dicedb/dice/testutils"
	"gotest.tools/v3/assert"
)

func TestRandomKey(t *testing.T) {
	conn := getLocalConnection()
	testCases := []struct {
		description string
		commands    []string
		expected    []interface{}
	}{
		{
			description: "no key returns (nil)",
			commands: []string{
				"RANDOMKEY",
			},
			expected: []interface{}{"(nil)"},
		},
		{
			description: "single defined key returns only that key",
			commands: []string{
				"SET name abc",
				"RANDOMKEY",
			},
			expected: []interface{}{
				"OK",
				"\"name\"",
			},
		},
		{
			description: "multiple defined keys return a random key from the defined key list",
			commands: []string{
				"SET name abc",
				"SET value def",
				"SET name2 ghi",
				"SET value-35 35",
				"RANDOMKEY",
				"RANDOMKEY",
				"RANDOMKEY",
			},
			expected: []interface{}{
				"OK",
				"OK",
				"OK",
				"OK",
				[]interface{}{"\"name\"", "\"value\"", "\"name2\"", "\"value-35\""},
				[]interface{}{"\"name\"", "\"value\"", "\"name2\"", "\"value-35\""},
				[]interface{}{"\"name\"", "\"value\"", "\"name2\"", "\"value-35\""},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			for i, cmd := range tc.commands {
				result := FireCommand(conn, cmd)

				// the result might be a single string, or a list of strings
				// for the 3nd TC, if the expected is a list of strings, we check
				// whether the result is one of the strings or not

				if arr, ok := tc.expected[i].([]interface{}); ok {
					assert.Assert(t, testutils.OneOf(result, arr))
				} else {
					assert.DeepEqual(t, tc.expected[i], result)
				}
			}
		})
	}
}
