package main

import "testing"

func Test_toSnakeCase(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s    string
		want string
	}{
		// Arrange
		{name: "Convert CamelCase to snake_case",
			s:    "CamelCase",
			want: "camel_case",
		},
		{
			name: "Convert PascalCase to snake_case",
			s:    "PascalCase",
			want: "pascal_case",
		},
		{
			name: "Convert mixedCase to snake_case",
			s:    "mixedCaseExample",
			want: "mixed_case_example",
		},
		{
			name: "Convert single word to snake_case",
			s:    "Word",
			want: "word",
		},
		{
			name: "Convert already snake_case to snake_case",
			s:    "already_snake_case",
			want: "already_snake_case",
		},
		{
			name: "Convert string with spaces to snake_case",
			s:    "string with spaces",
			want: "string_with_spaces",
		},
		{
			name: "Convert string with special characters to snake_case",
			s:    "special@char#string!",
			want: "special_char_string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSnakeCase(tt.s)
			if got != tt.want {
				t.Errorf("toSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
