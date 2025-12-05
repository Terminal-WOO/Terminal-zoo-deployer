.PHONY: build run dev stop clean test lint

# Docker image name and tag
IMAGE_NAME = ai-co
IMAGE_TAG = latest

# Docker compose commands
DC = docker compose

# Development commands
dev:
	npm install
	npm run dev

# Docker commands
build:
	$(DC) build

run:
	$(DC) up -d

stop:
	$(DC) down

# Clean up commands
clean:
	$(DC) down -v
	rm -rf node_modules
	rm -rf .nuxt
	rm -rf .output

# Application commands
install:
	npm ci

lint:
	npm run lint

test:
	npm run test

# Production commands
build-prod:
	docker build --platform linux/amd64 -t $(IMAGE_NAME):$(IMAGE_TAG) .

run-prod:
	docker run -d -p 3000:3000 $(IMAGE_NAME):$(IMAGE_TAG)

# Helper commands
logs:
	$(DC) logs -f

ps:
	$(DC) ps

restart:
	$(DC) restart

# Scaleway Container Registry
SCR_REGISTRY ?= rg.nl-ams.scw.cloud
SCR_NAMESPACE ?= nl-appstore-registry
FRONTEND_IMAGE = $(SCR_REGISTRY)/$(SCR_NAMESPACE)/ai-co
BACKEND_IMAGE = $(SCR_REGISTRY)/$(SCR_NAMESPACE)/ipo-deployer
IMAGE_TAG ?= latest

# Kubernetes namespace
K8S_NAMESPACE = nl-appstore-registry

# Build and push frontend image
build-frontend:
	docker build --platform linux/amd64 -t $(FRONTEND_IMAGE):$(IMAGE_TAG) .
	docker tag $(FRONTEND_IMAGE):$(IMAGE_TAG) $(FRONTEND_IMAGE):latest

push-frontend: build-frontend
	docker push $(FRONTEND_IMAGE):$(IMAGE_TAG)
	docker push $(FRONTEND_IMAGE):latest

# Build and push backend image
build-backend:
	docker build --platform linux/amd64 -t $(BACKEND_IMAGE):$(IMAGE_TAG) -f go/Dockerfile.multi-stage ./go
	docker tag $(BACKEND_IMAGE):$(IMAGE_TAG) $(BACKEND_IMAGE):latest

push-backend: build-backend
	docker push $(BACKEND_IMAGE):$(IMAGE_TAG)
	docker push $(BACKEND_IMAGE):latest

# Push both images
push-all: push-frontend push-backend

# Complete setup from .env file (pushes images and deploys everything)
setup-from-env:
	@echo "Running complete setup from .env file..."
	@bash deploy/setup-from-env.sh

# Kubernetes deployment commands
k8s-apply:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/frontend-deployment.yaml
	kubectl apply -f k8s/backend-deployment.yaml
	kubectl apply -f k8s/frontend-service.yaml
	kubectl apply -f k8s/backend-service.yaml
	kubectl apply -f k8s/ingress.yaml || echo "Ingress skipped (optional)"

# Create namespace only
k8s-create-namespace:
	kubectl apply -f k8s/namespace.yaml

# Apply HPA for cost-aware auto-scaling (optional)
k8s-apply-hpa:
	kubectl apply -f k8s/hpa.yaml

# Remove HPA to save costs (scale manually instead)
k8s-delete-hpa:
	kubectl delete -f k8s/hpa.yaml || true

k8s-delete:
	kubectl delete -f k8s/ingress.yaml || true
	kubectl delete -f k8s/frontend-service.yaml || true
	kubectl delete -f k8s/backend-service.yaml || true
	kubectl delete -f k8s/frontend-deployment.yaml || true
	kubectl delete -f k8s/backend-deployment.yaml || true
	kubectl delete -f k8s/configmap.yaml || true
	kubectl delete -f k8s/namespace.yaml || true

k8s-status:
	kubectl get pods -n $(K8S_NAMESPACE)
	kubectl get services -n $(K8S_NAMESPACE)
	kubectl get ingress -n $(K8S_NAMESPACE) || echo "No ingress configured"

# Get node IP (voor NodePort toegang)
k8s-get-node-ip:
	@echo "Node IP addresses:"
	@kubectl get nodes -o wide -n $(K8S_NAMESPACE) | grep -v NAME | awk '{print "  " $$7}'

# Get LoadBalancer IPs
k8s-get-lb-ips:
	@echo "LoadBalancer IPs:"
	@echo "Frontend:"
	@kubectl get svc nuxt-frontend-lb -n $(K8S_NAMESPACE) -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "  Not configured"
	@echo "Backend:"
	@kubectl get svc go-backend-lb -n $(K8S_NAMESPACE) -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "  Not configured"

# Apply NodePort services (voor gebruik zonder domain)
k8s-apply-nodeport:
	@echo "Creating namespace if it doesn't exist..."
	@kubectl apply -f k8s/namespace.yaml || true
	kubectl apply -f k8s/frontend-service-nodeport.yaml
	kubectl apply -f k8s/backend-service-nodeport.yaml
	@echo "NodePort services applied. Access via:"
	@echo "  Frontend: http://<node-ip>:30080"
	@echo "  Backend: http://<node-ip>:30088"
	@echo "Get node IP: make k8s-get-node-ip"

# Apply LoadBalancer services (voor gebruik zonder domain, maar met kosten)
k8s-apply-loadbalancer:
	@echo "Creating namespace if it doesn't exist..."
	@kubectl apply -f k8s/namespace.yaml || true
	kubectl apply -f k8s/frontend-service-loadbalancer.yaml
	kubectl apply -f k8s/backend-service-loadbalancer.yaml
	@echo "LoadBalancer services applied. Waiting for external IPs..."
	@echo "Check status: kubectl get svc -n $(K8S_NAMESPACE) -w"

k8s-logs-frontend:
	kubectl logs -f deployment/nuxt-frontend -n $(K8S_NAMESPACE)

k8s-logs-backend:
	kubectl logs -f deployment/go-backend -n $(K8S_NAMESPACE)

k8s-rollout-restart:
	kubectl rollout restart deployment/nuxt-frontend -n $(K8S_NAMESPACE)
	kubectl rollout restart deployment/go-backend -n $(K8S_NAMESPACE)

# Create Scaleway Container Registry secret
k8s-create-scr-secret:
	@echo "Creating Scaleway Container Registry secret..."
	@read -p "Enter Scaleway Secret Key: " secret_key; \
	kubectl create secret docker-registry scr-secret \
		--docker-server=$(SCR_REGISTRY) \
		--docker-username=nologin \
		--docker-password=$$secret_key \
		--namespace=$(K8S_NAMESPACE) \
		--dry-run=client -o yaml | kubectl apply -f -

# Help command
help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make dev          - Start development server"
	@echo "  make build        - Build using docker-compose"
	@echo "  make run          - Run using docker-compose"
	@echo "  make stop         - Stop all containers"
	@echo "  make clean        - Clean up containers and build files"
	@echo "  make install      - Install dependencies"
	@echo "  make lint         - Run linter"
	@echo "  make test         - Run tests"
	@echo ""
	@echo "Docker Production:"
	@echo "  make build-prod   - Build production Docker image"
	@echo "  make run-prod     - Run production Docker container"
	@echo ""
	@echo "Scaleway Container Registry:"
	@echo "  make build-frontend - Build frontend Docker image"
	@echo "  make push-frontend  - Build and push frontend image to SCR"
	@echo "  make build-backend  - Build backend Docker image"
	@echo "  make push-backend   - Build and push backend image to SCR"
	@echo "  make push-all       - Build and push both images to SCR"
	@echo ""
	@echo "Complete Setup (from .env):"
	@echo "  make setup-from-env - Complete setup: push images + deploy (gebruikt .env)"
	@echo ""
	@echo "Kubernetes Deployment:"
	@echo "  make k8s-apply              - Apply all Kubernetes manifests"
	@echo "  make k8s-delete             - Delete all Kubernetes resources"
	@echo "  make k8s-status             - Show Kubernetes deployment status"
	@echo "  make k8s-logs-frontend      - Show frontend logs"
	@echo "  make k8s-logs-backend       - Show backend logs"
	@echo "  make k8s-rollout-restart    - Restart deployments"
	@echo "  make k8s-create-scr-secret  - Create SCR image pull secret"
	@echo "  make k8s-apply-hpa          - Apply HPA for auto-scaling (optional)"
	@echo "  make k8s-delete-hpa         - Remove HPA to save costs"
	@echo "  make k8s-get-node-ip        - Get node IP addresses"
	@echo "  make k8s-get-lb-ips         - Get LoadBalancer IP addresses"
	@echo "  make k8s-create-namespace    - Create namespace only"
	@echo "  make k8s-apply-nodeport      - Apply NodePort services (no domain needed)"
	@echo "  make k8s-apply-loadbalancer  - Apply LoadBalancer services (costs ~€20/mo)"
	@echo ""
	@echo "Helper commands:"
	@echo "  make logs         - View container logs"
	@echo "  make ps           - List running containers"
	@echo "  make restart      - Restart containers"