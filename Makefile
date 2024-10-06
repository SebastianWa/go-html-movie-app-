run:
	templ generate
	go run .

tw:
	npx tailwindcss -i ./static/css/main.css -o ./static/css/output.css --watch 

templ:
	templ generate