# go-wkhtml-container-image-cloud-run
This repository contains the Dockerfile for the go-wkhtmltopdf container image for google cloud run or aws ecs.

# scripts

Makefile contains the following scripts:

## build
```bash
docker build -t go-html-to-pdf .
```

## start test server
```bash
docker run --rm -p 8080:8080 go-html-to-pdf
```