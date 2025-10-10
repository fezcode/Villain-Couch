package error

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateErrors() error {
	err1 := errors.New("password is too short")
	err2 := errors.New("username contains invalid characters")
	var err3 error // A nil error to ensure it gets filtered out.

	return errors.Join(err1, err2, err3)
}

// TestValidateAndUnwrapErrors tests that our validate function works as expected
// and that we can correctly iterate through the joined errors.
func TestValidateAndUnwrapErrors(t *testing.T) {
	expectedErrors := []string{
		"password is too short",
		"username contains invalid characters",
	}
	err := generateErrors()

	// The error should not be nil.
	if err == nil {
		t.Fatal("validate() returned nil, but an error was expected")
	}

	count := UnwrapErrors(err)
	assert.Equal(t, count, len(expectedErrors))
}
