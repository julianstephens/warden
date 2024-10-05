package warden

func DefaultIfNil[T any](value interface{}, defaultValue interface{}) T {
	if value != nil {
		return value.(T)
	}
	return defaultValue.(T)
}
