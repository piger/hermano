package config

import (
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	UserKey  string   `toml:"user_key"`
	APIToken string   `toml:"api_token"`
	Ignored  []string `toml:"ignored"`
}

func ReadConfig(filename string) (*Config, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	contents, err := io.ReadAll(fh)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := toml.Unmarshal(contents, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
