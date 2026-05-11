package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"library-management-system/backend/controllers"
	"library-management-system/backend/initializers"
	"library-management-system/backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

var runEnv string

func init() {
	envFlag := flag.String("env", "", "required: development or production")
	flag.Parse()

	runEnv = strings.TrimSpace(*envFlag)
	if runEnv == "" {
		log.Fatal("missing -env: use -env=development or -env=production (run from the backend directory so .env.<env> is found)")
	}

	if runEnv != initializers.EnvDevelopment && runEnv != initializers.EnvProduction {
		log.Fatalf("invalid -env %q: use %q or %q", runEnv, initializers.EnvDevelopment, initializers.EnvProduction)
	}

	if err := initializers.LoadEnv(runEnv); err != nil {
		log.Fatalf("env: %v", err)
	}
	if runEnv == initializers.EnvProduction && strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS")) == "" {
		log.Fatal("CORS_ALLOW_ORIGINS must be set in .env.production (comma-separated origins, e.g. https://app.example.com)")
	}
	if err := initializers.ConnectDB(runEnv); err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := initializers.SyncDB(); err != nil {
		log.Fatalf("database sync: %v", err)
	}
}

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "Library API",
		ServerHeader: "Fiber",
	})

	for _, mw := range middlewares.Stack(runEnv) {
		app.Use(mw)
	}

	app.Get("/health", controllers.Health)

	// book routes
	app.Post("/book", controllers.CreateBook)
	app.Delete("/book/:id", controllers.DeleteBook)

	addr := strings.TrimSpace(os.Getenv("SERVER_ADDR"))
	if addr == "" {
		if runEnv == initializers.EnvProduction {
			addr = ":8080"
		} else {
			addr = ":3000"
		}
	}

	log.Printf("environment=%s listening on %s", runEnv, addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
