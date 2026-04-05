//go:build !windows && !darwin
package resolver

func findVlcOnWindows() (string, bool, error) { return "", false, nil }
func findVlcOnDarwin() (string, bool, error) { return "", false, nil }
