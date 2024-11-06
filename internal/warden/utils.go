package warden

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

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

	err = json.Unmarshal(b, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf("\n"+format+"\n", a...)
}

func PPrint(data interface{}) (n int, err error) {
	transformer := text.NewJSONTransformer("", strings.Repeat(" ", 2))
	pJson := transformer(data)
	return fmt.Println(pJson)
}

func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModeDir|0755)
	}
	return nil
}

func LoadJSON[T any](filename string) (T, error) {
	var data T
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}
	return data, json.Unmarshal(fileData, &data)
}
