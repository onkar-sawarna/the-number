TEMPL ?= go run github.com/a-h/templ/cmd/templ@v0.2.793

.PHONY: generate copy css og vendor run test tidy install

copy:
	python3 scripts/gen_copy.py

generate: copy
	$(TEMPL) generate

css:
	npx tailwindcss -i web/static/css/input.css -o web/static/css/app.css --minify

og:
	python3 scripts/gen_og.py

vendor:
	python3 scripts/vendor.py

run: generate
	PORT=47321 go run ./cmd/server

test:
	go test ./internal/calc ./internal/ai ./internal/handlers ./web/templates

tidy:
	go mod tidy

install:
	go install github.com/a-h/templ/cmd/templ@v0.2.793
	go install github.com/air-verse/air@latest
