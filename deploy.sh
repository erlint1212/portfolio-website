#!/bin/bash
# deploy.sh — Build locally, push to GHCR, let ArgoCD handle the rest
set -e

IMAGE="ghcr.io/erlint1212/portfolio-website"
TAG=$(git rev-parse --short HEAD)

echo "Building image: $IMAGE:$TAG"
docker build -t "$IMAGE:$TAG" -t "$IMAGE:latest" .

echo "Pushing to GHCR"
docker push "$IMAGE:$TAG"
docker push "$IMAGE:latest"

echo "Updating deployment manifest"
sed -i "s|image: ghcr.io/erlint1212/portfolio-website:.*|image: ghcr.io/erlint1212/portfolio-website:$TAG|" kubernetes/portfolio-deployment.yaml

echo "Committing and pushing"
git add kubernetes/portfolio-deployment.yaml
git diff --staged --quiet || git commit -m "deploy: update image to $TAG"
git push

echo "Done! ArgoCD will pick it up."
