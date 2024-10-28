package common

type Params interface{}

type LocalStorageParams struct {
	Params
	Location string
}
