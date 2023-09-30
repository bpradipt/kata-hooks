package internal

import (
	"testing"
)

func TestConvertOptionsToString(t *testing.T) {
	testCases := []struct {
		name     string
		options  []string
		expected string
	}{
		{
			name:     "empty options",
			options:  []string{},
			expected: "",
		},
		{
			name:     "single option",
			options:  []string{"foo"},
			expected: "foo",
		},
		{
			name:     "multiple options",
			options:  []string{"foo", "bar", "baz"},
			expected: "foo,bar,baz",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ConvertOptionsToString(tc.options)
			if actual != tc.expected {
				t.Errorf("expected %q, but got %q", tc.expected, actual)
			}
		})
	}
}
