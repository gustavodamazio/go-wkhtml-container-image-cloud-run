#!/bin/bash

# Set default variables
PROJECT_ID="<project-id>"
REPOSITORY_NAME="docker-images"
IMAGE_NAME="go-html-to-pdf"
TAG="v1"
REGION="us-central1"
SERVICE_NAME="html-to-pdf"

# Example of environment variables
ENV_VARS="GCS_BUCKET_NAME=<project-id>-relatorios-tmp,ENVIRONMENT=production"

# Function to display help
usage() {
    echo "Usage: $0 [-p project_id] [-r repository_name] [-i image_name] [-t tag] [-g region] [-s service_name] [-e env_vars]"
    exit 1
}

# Parse arguments
while getopts "p:r:i:t:g:s:e:" opt; do
    case $opt in
        p) PROJECT_ID="$OPTARG" ;;
        r) REPOSITORY_NAME="$OPTARG" ;;
        i) IMAGE_NAME="$OPTARG" ;;
        t) TAG="$OPTARG" ;;
        g) REGION="$OPTARG" ;;
        s) SERVICE_NAME="$OPTARG" ;;
        e) ENV_VARS="$OPTARG" ;;  # Capture environment variables
        *) usage ;;
    esac
done

# Build the Docker image
docker build --platform linux/amd64 -t "${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY_NAME}/${IMAGE_NAME}:${TAG}" .

# Authenticate with Artifact Registry
gcloud auth configure-docker "${REGION}-docker.pkg.dev"

# Push the image to Artifact Registry
docker push "${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY_NAME}/${IMAGE_NAME}:${TAG}"

# Deploy the image to Cloud Run with environment variables
gcloud run deploy "${SERVICE_NAME}" \
    --image "${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY_NAME}/${IMAGE_NAME}:${TAG}" \
    --region "${REGION}" \
    --platform managed \
    --allow-unauthenticated \
    --memory "2Gi" \
    --concurrency "5" \
    --timeout "300" \
    --set-env-vars "${ENV_VARS}"
