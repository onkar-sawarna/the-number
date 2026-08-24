TEMPL ?= go run github.com/a-h/templ/cmd/templ@v0.2.793

.PHONY: generate css run test tidy install

generate:
	$(TEMPL) generate

css:
	npx tailwindcss -i web/static/css/input.css -o web/static/css/app.css --minify

run: generate
	PORT=47321 go run ./cmd/server

test:
	go test ./internal/calc ./internal/ai

tidy:
	go mod tidy

install:
	go install github.com/a-h/templ/cmd/templ@v0.2.793
	go install github.com/air-verse/air@latest
