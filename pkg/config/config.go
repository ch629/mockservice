package config

type Config struct {
	API API `mapstructure:"api"`
}

type API struct {
	Port int `mapstructure:"port"`
}
