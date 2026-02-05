package config

type Config struct {
	Port   string `json:"port"`
	Remote string `json:"remote"`
}

func LoadConfig() *Config {
	return &Config{}
}
