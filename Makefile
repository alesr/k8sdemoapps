.DEFAULT_GOAL := help

PROJECT_NAME := Demo App Deployment

.PHONY: help
help:
	@echo "------------------------------------------------------------------------"
	@echo "${PROJECT_NAME}"
	@echo "------------------------------------------------------------------------"
	@grep -E '^[a-zA-Z0-9_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

APP_NAME ?= demoapp1
APP_PATH ?= ../cmd/$(APP_NAME)

.PHONY: build
build:  ## Build the Docker image for the app
	@echo "Building Docker image for application: $(APP_NAME)"
	docker build --build-arg APP_NAME=$(APP_NAME) --build-arg APP_PATH=$(APP_PATH) -t $(APP_NAME) -f build/Dockerfile .

.PHONY: build-all
build-all:  ## Build Docker images for all demo apps
	@$(MAKE) build APP_NAME=demoapp1
	@$(MAKE) build APP_NAME=demoapp2
	@$(MAKE) build APP_NAME=demoapp3

.PHONY: up
up:  ## Run the application with Docker Compose
	@if [ -z "$(shell docker images -q $(APP_NAME))" ]; then \
		echo "Image for $(APP_NAME) not found, building..."; \
		$(MAKE) build; \
	fi; \
	@echo "Starting application: $(APP_NAME)"
	docker compose -f build/docker-compose.yml up $(APP_NAME) -d --build

.PHONY: up-all
up-all:  ## Run all applications with Docker Compose
	@echo "Starting all applications"
	docker compose -f build/docker-compose.yml up -d --build

.PHONY: down
down:  ## Stop all running Docker containers
	@echo "Stopping all applications"
	docker compose -f build/docker-compose.yml down

.PHONY: clean
clean:  ## Clean up Docker containers and images
	@echo "Cleaning up Docker containers and images"
	docker compose -f build/docker-compose.yml down --rmi all --volumes --remove-orphans
	docker system prune -f

NAMESPACE_DEMOAPPS := k8sdemoapps
NAMESPACE_TRAEFIK := traefik

.PHONY: k8s-setup
k8s-setup:  ## Set up Kubernetes namespaces and install Traefik
	@echo "Setting up Kubernetes namespaces..."
	kubectl delete namespace $(NAMESPACE_TRAEFIK) || true
	kubectl create namespace $(NAMESPACE_TRAEFIK)
	kubectl delete namespace $(NAMESPACE_DEMOAPPS) || true
	kubectl create namespace $(NAMESPACE_DEMOAPPS)
	@echo "Namespaces set up: $(NAMESPACE_DEMOAPPS), $(NAMESPACE_TRAEFIK)"

.PHONY: k8s-install-traefik
k8s-install-traefik: k8s-setup  ## Install Traefik using Helm in Kubernetes
	@echo "Installing Traefik..."
	helm repo add traefik https://helm.traefik.io/traefik
	helm repo update
	helm install traefik traefik/traefik --namespace $(NAMESPACE_TRAEFIK) -f build/traefik/config.yaml
	@echo "Traefik installed."

.PHONY: k8s-deployment-apply
k8s-deployment-apply:  ## Apply Kubernetes deployment manifests for demo apps
	kubectl apply -f build/k8s-manifests/demoapp1-deployment.yaml
	kubectl apply -f build/k8s-manifests/demoapp2-deployment.yaml
	kubectl apply -f build/k8s-manifests/demoapp3-deployment.yaml

.PHONY: k8s-deployment-delete
k8s-deployment-delete:  ## Delete demo app deployments in Kubernetes
	kubectl delete deployment demoapp1 demoapp2 demoapp3 -n $(NAMESPACE_DEMOAPPS)

.PHONY: k8s-port-forward-traefik
k8s-port-forward-traefik:  ## Port-forward Traefik dashboard in Kubernetes
	@echo "Port-forwarding Traefik dashboard on 'localhost:8080/dashboard/'"
	kubectl port-forward -n $(NAMESPACE_TRAEFIK) $$(kubectl get pods -n $(NAMESPACE_TRAEFIK) -l app.kubernetes.io/name=traefik -o jsonpath='{.items[0].metadata.name}') 8080:9000
