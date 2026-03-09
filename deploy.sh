#!/usr/bin/env bash
# Simple deployment before proper CI/CD with github actions or something simmilar
set -e

SERVER="thinkpad"
IMAGE="portfolio:latest"
HASH_FILE=".last-deploy-hash"

echo "==> Building Docker image..."
docker build -t $IMAGE .

# Get the new image ID
NEW_HASH=$(docker inspect --format='{{.Id}}' $IMAGE)

# Check if the image actually changed
OLD_HASH=""
if [ -f "$HASH_FILE" ]; then
    OLD_HASH=$(cat "$HASH_FILE")
fi

if [ "$NEW_HASH" != "$OLD_HASH" ]; then
    echo "==> Image changed, transferring..."
    docker save $IMAGE | gzip > /tmp/portfolio-image.tar.gz
    scp /tmp/portfolio-image.tar.gz $SERVER:~/portfolio-image.tar.gz
    ssh -t $SERVER "gunzip -c ~/portfolio-image.tar.gz | sudo k3s ctr images import - && rm ~/portfolio-image.tar.gz"
    echo "$NEW_HASH" > "$HASH_FILE"
    echo "==> Restarting deployment..."
    ssh $SERVER "kubectl rollout restart deployment/portfolio-web"
else
    echo "==> Image unchanged, skipping transfer."
fi

# Only sync manifests that changed
echo "==> Syncing Kubernetes manifests..."
rsync -avz --checksum kubernetes/ $SERVER:~/kubernetes/

echo "==> Applying manifests..."
ssh $SERVER "kubectl apply -f ~/kubernetes/"

echo "==> Status:"
ssh $SERVER "kubectl get pods"
