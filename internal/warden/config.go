package warden

type Config struct {
	ID     string         `json:"id"`
	Params map[string]int `json:"params"`
}

func CreateConfig(params map[string]int) (Config, error) {
	var conf Config

	conf.ID = NewID().String()
	conf.Params = params

	return conf, nil
}
