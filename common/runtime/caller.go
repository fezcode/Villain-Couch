package runtime

import (
	"runtime"
)

// GetCallerGrandparent retrieves the name of the function that is two levels
// up the call stack. For example, if A() calls B(), and B() calls this function,
// it will return the name of A().
// It returns the fully qualified function name (e.g., "main.main") and a boolean
// indicating if the call was successful.
func GetCallerGrandparent() (name string, ok bool) {
	// We use a skip value of 3 to get the grandparent caller.
	// 0: runtime.Callers itself
	// 1: GetCallerGrandparent (this function)
	// 2: The immediate caller (the "parent")
	// 3: The caller of the caller (the "grandparent")
	const skipFrames = 3

	// Get the program counter for the target frame.
	pc := make([]uintptr, 1)
	n := runtime.Callers(skipFrames, pc)
	if n == 0 {
		// Not enough callers in the stack.
		return "", false
	}

	// Get the function object from the program counter.
	fn := runtime.FuncForPC(pc[0])
	if fn == nil {
		return "", false
	}

	// The full name includes the package path. We can optionally trim it.
	fullName := fn.Name()
	// Example of trimming the package path for a cleaner name:
	// lastSlash := strings.LastIndex(fullName, "/")
	// if lastSlash >= 0 {
	// 	 fullName = fullName[lastSlash+1:]
	// }

	return fullName, true
}
