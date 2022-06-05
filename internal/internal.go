package internal

import "fmt"

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
