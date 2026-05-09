package initializers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Environment names must match the suffix of .env.<name> files.
const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

// LoadEnv loads `.env.<environment>` from the current working directory.
func LoadEnv(environment string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	filename := filepath.Join(wd, fmt.Sprintf(".env.%s", environment))
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("load %s: %w", filename, err)
	}

	log.Printf("loaded environment file: %s", filename)
	return nil
}
