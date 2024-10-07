#!/bin/bash

# Set default variables
PROJECT_ID="your-project-id"
REPOSITORY_NAME="my-repository"
REGION="us-central1"

# Function to display help
usage() {
    echo "Usage: $0 [-p project_id] [-r repository_name] [-g region]"
    exit 1
}

# Parse arguments
while getopts "p:r:g:" opt; do
    case $opt in
        p) PROJECT_ID="$OPTARG" ;;
        r) REPOSITORY_NAME="$OPTARG" ;;
        g) REGION="$OPTARG" ;;
        *) usage ;;
    esac
done

# Authenticate with Google Cloud
echo "Authenticating with Google Cloud..."
gcloud auth login || { echo "Error: Failed to authenticate with Google Cloud"; exit 1; }

# Configure the project
echo "Setting project $PROJECT_ID..."
gcloud config set project "$PROJECT_ID" || { echo "Error: Failed to set the project"; exit 1; }

# Enable necessary APIs
echo "Enabling necessary APIs..."
gcloud services enable artifactregistry.googleapis.com \
    && gcloud services enable run.googleapis.com \
    || { echo "Error: Failed to enable APIs"; exit 1; }

# Create a repository in Artifact Registry
echo "Creating repository $REPOSITORY_NAME in Artifact Registry in region $REGION..."
gcloud artifacts repositories create "$REPOSITORY_NAME" \
    --repository-format=docker \
    --location="$REGION" \
    --description="Docker repository" \
    || { echo "Error: Failed to create the repository"; exit 1; }

echo "Environment successfully configured!"
