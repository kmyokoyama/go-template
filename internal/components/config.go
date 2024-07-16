package components

import (
	"os"

	"github.com/joho/godotenv"
)

type Config interface {
	Get(string) string
}

type config struct {}

func (c *config) Get(key string) string {
	return os.Getenv(key)
}

func NewConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic("unable to load env vars, panicking")
	}
	
	return &config{}
}

// Verify interface compliance.
var _ Config = (*config)(nil)