package optional

func FirstOrEmpty[T any](value []T) T {
	if len(value) > 0 {
		return value[0]
	}
	var empty T
	return empty
}
