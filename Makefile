watch-css:
	npx tailwindcss -i ./styles/input.css -o ./styles/output.css --watch

serve:
	go run .

.PHONY: watch build
