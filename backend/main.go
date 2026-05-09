package main

import (
	"flag"
	"log"
	"os"

	"library-management-system/backend/controllers"
	"library-management-system/backend/initializers"
	"library-management-system/backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

var runEnv string

func init() {
	envFlag := flag.String("env", "", "environment: development or production (default: APP_ENV or development)")
	flag.Parse()

	runEnv = *envFlag
	if runEnv == "" {
		runEnv = os.Getenv("APP_ENV")
	}
	if runEnv == "" {
		runEnv = initializers.EnvDevelopment
	}

	if runEnv != initializers.EnvDevelopment && runEnv != initializers.EnvProduction {
		log.Fatalf("invalid -env / APP_ENV %q: use %q or %q", runEnv, initializers.EnvDevelopment, initializers.EnvProduction)
	}

	if err := initializers.LoadEnv(runEnv); err != nil {
		log.Fatalf("env: %v", err)
	}
	if err := initializers.ConnectDB(runEnv); err != nil {
		log.Fatalf("database: %v", err)
	}
}

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "Library API",
		ServerHeader: "Fiber",
	})

	for _, mw := range middlewares.Stack() {
		app.Use(mw)
	}

	app.Get("/health", controllers.Health)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":3000"
	}

	log.Printf("environment=%s listening on %s", runEnv, addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
