package warden

import "fmt"

const (
	ExitCodeErr       = 1
	ExitCodeInterrupt = 2
)

type InvalidStoreError struct {
	Msg string
}

func (error InvalidStoreError) Error() string {
	return fmt.Sprintf("invalid store: %s", error.Msg)
}

type InvalidPasswordError struct {
	Msg string
}

func (error InvalidPasswordError) Error() string {
	return fmt.Sprintf("invalid password: %s", error.Msg)
}

type InvalidArgumentError struct {
	Expecting string
	Got       string
}

func (error InvalidArgumentError) Error() string {
	return fmt.Sprintf("expecting %s, got %s", error.Expecting, error.Got)
}
