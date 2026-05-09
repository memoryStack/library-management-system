package initializers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Environment names must match suffix of .env.<name> files.
const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

// LoadEnv loads `.env.<environment>` from the process working directory (or BACKEND_ROOT if set).
func LoadEnv(environment string) error {
	root := os.Getenv("BACKEND_ROOT")
	if root == "" {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}
		root = wd
	}

	filename := filepath.Join(root, fmt.Sprintf(".env.%s", environment))
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("load %s: %w", filename, err)
	}

	log.Printf("loaded environment file: %s", filename)
	return nil
}
