.PHONY: dev tidy run-dev run-prod

# Install CompileDaemon once: go install github.com/githubnemo/CompileDaemon@latest
dev:
	CompileDaemon -build="go build -o tmp/main ." -command="./tmp/main -env=development" -directory=. -pattern="(.+\.go|.+\.env.*)$$" -recursive=true

tidy:
	go mod tidy

run-dev:
	go run . -env=development

run-prod:
	go run . -env=production
