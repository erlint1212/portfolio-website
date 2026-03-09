#!/usr/bin/env bash
# Simple deployment before proper CI/CD with github actions or something simmilar
set -e

SERVER="thinkpad"
IMAGE="portfolio:latest"

echo "==> Building Docker image..."
docker build -t $IMAGE .

echo "==> Saving and transferring image..."
docker save $IMAGE | gzip | ssh $SERVER "gunzip | sudo k3s ctr images import -"

echo "==> Applying Kubernetes manifests..."
scp -r kubernetes/ $SERVER:~/server/kubernetes/
ssh $SERVER "kubectl apply -f ~/server/kubernetes/"

echo "==> Restarting deployment..."
ssh $SERVER "kubectl rollout restart deployment/portfolio"

echo "==> Done! Checking pods..."
ssh $SERVER "kubectl get pods"
