build:
	docker build -t go-html-to-pdf .
run:
	docker run --rm -p 8080:8080 go-html-to-pdf
setup-dev:
	gcloud config set project <project-id> && ./scripts/setup_environment.sh -p <project-id> -r docker-images -g us-central1
deploy-dev:
	gcloud config set project <project-id> && ./scripts/deploy_cloud_run.sh -p <project-id> -r docker-images -i go-html-to-pdf -t v1 -g us-central1 -s html-to-pdf -e "GCS_BUCKET_NAME=<project-id>-relatorios-tmp,ENVIRONMENT=production"