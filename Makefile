build:
	docker build -t go-html-to-pdf .
run:
	docker run --rm -p 8080:8080 go-html-to-pdf