package error

import (
	"fmt"
)

// UnwrapErrors checks if an error is a joined error (from errors.Join)
// or a single error. It prints each individual error message and returns
// the total count of errors found.
func UnwrapErrors(err error) int {
	// First, handle the case of no error.
	if err == nil {
		return 0
	}

	// Attempt to unwrap the error into a slice of errors.
	// This is the standard way to check for errors created by errors.Join.
	if unwrap, ok := err.(interface {
		Unwrap() []error
	}); ok {
		// It's a joined error. Loop through the slice of original errors.
		unwrappedErrors := unwrap.Unwrap()
		for i, individualErr := range unwrappedErrors {
			fmt.Printf("  - Error %d: %s\n", i+1, individualErr.Error())
		}
		return len(unwrappedErrors)
	}

	// If the type assertion fails, it's a single error.
	// Print the single error message and return a count of 1.
	fmt.Printf("  - Error: %s\n", err.Error())
	return 1
}
