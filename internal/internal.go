package internal

import (
	"fmt"
	"strings"
)

// resolve returns the dereferenced pointer value if the pointer is non-nil,
// otherwise it returns the given default.
func resolve[T any](v *T, def T) T {
	if v != nil {
		return *v
	}
	return def
}

// failF creates a function for wrapping other error functions.
func failF(funcName string) func(error) error {
	return func(err error) error {
		return fmt.Errorf("%s: %w", funcName, err)
	}
}

// nullIf returns nil if the given string is empty or only contains whitespaces, otherwise
// it returns a pointer to the string value.
func nullIf[T any](value T, pred func(T) bool) *T {
	if pred(value) {
		return nil
	}
	return &value
}

func textEmpty(v string) bool {
	return v == "" || strings.TrimSpace(v) == ""
}
