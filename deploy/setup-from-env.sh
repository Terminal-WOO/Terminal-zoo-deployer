#!/bin/bash

# Complete Setup Script - Gebruikt secrets uit .env bestand
# Dit script doet alles: images pushen, secrets aanmaken, deployments deployen

set -e  # Stop bij errors

# Load .env file
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "❌ .env bestand niet gevonden!"
    exit 1
fi

# Check required variables
if [ -z "$SCW_ACCESS_KEY" ] || [ -z "$SCW_SECRET_KEY" ]; then
    echo "❌ SCW_ACCESS_KEY of SCW_SECRET_KEY niet gevonden in .env!"
    exit 1
fi

echo "✅ .env bestand geladen"
echo ""

# Configuration
SCR_REGISTRY="rg.nl-ams.scw.cloud"
SCR_NAMESPACE="nl-appstore-registry"
K8S_NAMESPACE="nl-appstore-registry"
FRONTEND_IMAGE="$SCR_REGISTRY/$SCR_NAMESPACE/ai-co"
BACKEND_IMAGE="$SCR_REGISTRY/$SCR_NAMESPACE/ai-co"

# 1. Docker login
echo "=== Stap 1: Docker Login ==="
docker login $SCR_REGISTRY -u nologin -p "$SCW_SECRET_KEY" || {
    echo "❌ Docker login gefaald!"
    exit 1
}
echo "✅ Docker ingelogd"
echo ""

# 2. Build en push images
echo "=== Stap 2: Building and Pushing Images ==="
echo "Building frontend..."
docker build --platform linux/amd64 -t $FRONTEND_IMAGE:latest . || {
    echo "❌ Frontend build gefaald!"
    exit 1
}
docker tag $FRONTEND_IMAGE:latest $FRONTEND_IMAGE:latest

echo "Pushing frontend..."
docker push $FRONTEND_IMAGE:latest || {
    echo "❌ Frontend push gefaald!"
    exit 1
}
echo "✅ Frontend image gepusht"

echo "Building backend..."
docker build --platform linux/amd64 -t $BACKEND_IMAGE:latest -f go/Dockerfile.multi-stage ./go || {
    echo "❌ Backend Docker build gefaald!"
    exit 1
}
docker tag $BACKEND_IMAGE:latest $BACKEND_IMAGE:latest

echo "Pushing backend..."
docker push $BACKEND_IMAGE:latest || {
    echo "❌ Backend push gefaald!"
    exit 1
}
echo "✅ Backend image gepusht"
echo ""

# 3. ConfigMap
echo "=== Stap 3: Creating ConfigMap ==="
kubectl apply -f k8s/configmap.yaml || {
    echo "❌ ConfigMap aanmaken gefaald!"
    exit 1
}
echo "✅ ConfigMap aangemaakt"
echo ""

# 4. Secrets
echo "=== Stap 4: Creating Secrets ==="

# App secrets (gebruik lege token als niet gezet)
API_TOKEN="${NUXT_EXTERNAL_API_AUTH:-Bearer }"
kubectl create secret generic app-secrets \
  --from-literal=NUXT_EXTERNAL_API_AUTH="$API_TOKEN" \
  -n $K8S_NAMESPACE --dry-run=client -o yaml | kubectl apply -f - || {
    echo "⚠️ App secrets al bestaat, overschrijven..."
    kubectl delete secret app-secrets -n $K8S_NAMESPACE --ignore-not-found=true
    kubectl create secret generic app-secrets \
      --from-literal=NUXT_EXTERNAL_API_AUTH="$API_TOKEN" \
      -n $K8S_NAMESPACE
}
echo "✅ App secrets aangemaakt"

# Registry secret
kubectl create secret docker-registry scr-secret \
  --docker-server=$SCR_REGISTRY \
  --docker-username=nologin \
  --docker-password="$SCW_SECRET_KEY" \
  -n $K8S_NAMESPACE --dry-run=client -o yaml | kubectl apply -f - || {
    echo "⚠️ Registry secret al bestaat, overschrijven..."
    kubectl delete secret scr-secret -n $K8S_NAMESPACE --ignore-not-found=true
    kubectl create secret docker-registry scr-secret \
      --docker-server=$SCR_REGISTRY \
      --docker-username=nologin \
      --docker-password="$SCW_SECRET_KEY" \
      -n $K8S_NAMESPACE
}
echo "✅ Registry secret aangemaakt"
echo ""

# 5. Deployments
echo "=== Stap 5: Deploying Applications ==="
kubectl apply -f k8s/frontend-deployment.yaml || {
    echo "❌ Frontend deployment gefaald!"
    exit 1
}
kubectl apply -f k8s/backend-deployment.yaml || {
    echo "❌ Backend deployment gefaald!"
    exit 1
}
echo "✅ Deployments aangemaakt"
echo ""

# 6. Wacht op pods
echo "=== Stap 6: Waiting for Pods ==="
echo "Wachten op frontend pod..."
kubectl wait --for=condition=ready pod -l app=nuxt-frontend -n $K8S_NAMESPACE --timeout=5m || {
    echo "⚠️ Frontend pod start niet binnen timeout"
    kubectl get pods -l app=nuxt-frontend -n $K8S_NAMESPACE
}

echo "Wachten op backend pod..."
kubectl wait --for=condition=ready pod -l app=go-backend -n $K8S_NAMESPACE --timeout=5m || {
    echo "⚠️ Backend pod start niet binnen timeout"
    kubectl get pods -l app=go-backend -n $K8S_NAMESPACE
}
echo "✅ Pods zijn ready"
echo ""

# 7. NodePort services
echo "=== Stap 7: Deploying NodePort Services ==="
kubectl apply -f k8s/frontend-service-nodeport.yaml || {
    echo "❌ Frontend NodePort service gefaald!"
    exit 1
}
kubectl apply -f k8s/backend-service-nodeport.yaml || {
    echo "❌ Backend NodePort service gefaald!"
    exit 1
}
echo "✅ NodePort services aangemaakt"
echo ""

# 8. Get node IP
echo "=== Stap 8: Node IP Address ==="
NODE_IP=$(kubectl get nodes -o wide | grep -v NAME | awk '{print $7}' | head -1)
if [ -z "$NODE_IP" ]; then
    NODE_IP=$(kubectl get nodes -o wide | grep -v NAME | awk '{print $6}' | head -1)
fi

if [ -z "$NODE_IP" ]; then
    echo "⚠️ Kon node IP niet vinden"
else
    echo "✅ Node IP gevonden: $NODE_IP"
    echo ""
    echo "🌐 Toegang:"
    echo "   Frontend: http://$NODE_IP:30080"
    echo "   Backend: http://$NODE_IP:30088"
    echo ""
fi

# 9. Status
echo "=== Stap 9: Final Status ==="
kubectl get pods -n $K8S_NAMESPACE
echo ""
kubectl get svc -n $K8S_NAMESPACE
echo ""

# 10. Test (als node IP bekend is)
if [ ! -z "$NODE_IP" ]; then
    echo "=== Stap 10: Testing ==="
    echo "Testing frontend health..."
    curl -s http://$NODE_IP:30080/api/health || echo "⚠️ Frontend health check gefaald"
    echo ""
    echo "Testing backend health..."
    curl -s http://$NODE_IP:30088/health || echo "⚠️ Backend health check gefaald"
    echo ""
fi

echo "✅ Setup compleet!"
echo ""
echo "📝 Volgende stappen:"
echo "   1. Configureer firewall rules voor poorten 30080 en 30088"
echo "   2. Test toegang: http://$NODE_IP:30080"
echo "   3. Check logs: make k8s-logs-frontend"
echo ""

