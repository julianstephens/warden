package config

type Config struct {
	ID string `json:"id"`
}

func CreateConfig() (Config, error) {
	var conf Config

	conf.ID = NewID().String()

	return conf, nil
}
