package helpers

func ValueOrDefault[T any](val *T, defaultValue T) T {
	if val == nil {
		return defaultValue
	}

	return *val
}
