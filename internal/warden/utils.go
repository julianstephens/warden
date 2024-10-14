package warden

import "encoding/json"

func DefaultIfNil[T any](value interface{}, defaultValue interface{}) T {
	if value != nil {
		return value.(T)
	}
	return defaultValue.(T)
}

func MapToStruct[T any](value any) (T, error) {
	var res T

	b, err := json.Marshal(value)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return res, err
	}

	return res, nil
}
