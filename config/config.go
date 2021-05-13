package config

import (
	"encoding/json"
	"os"

	"github.com/juju/errors"
)

type Config interface {
	Values() *Values
}

type FileConfig struct {
	values *Values
}

// New loads the configuration values for the app
func NewFileConfig() (*FileConfig, error) {
	file, err := os.ReadFile("./config/config.json")
	if err != nil {
		return nil, errors.NotFoundf(err.Error())
	}

	var config *Values
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, errors.NotValidf(err.Error())
	}

	return &FileConfig{config}, nil
}

// Values return the app configuration values
func (c *FileConfig) Values() *Values {
	return c.values
}
