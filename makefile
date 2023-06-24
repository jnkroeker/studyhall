SHELL := /bin/bash

# ======================================================================
# Fresh Machine Install

GOLANG          := golang:1.20
ALPINE          := alpine:3.18
KIND            := kindest/node:v1.27.3
POSTGRES        := postgres:15.3
VAULT           := hashicorp/vault:1.13
GRAFANA         := grafana/grafana:9.5.3
PROMETHEUS      := prom/prometheus:v2.44.0
TEMPO           := grafana/tempo:2.1.1
TELEPRESENCE    := datawire/ambassador-telepresence-manager:2.14.0

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-brew-common:
	brew update
	brew tap hashicorp/tap
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli
	brew list vault || brew install vault

dev-brew: dev-brew-common
	brew list datawire/blackbird/telepresence || brew install datawire/blackbird/telepresence

dev-docker:
	docker pull $(GOLANG)
	docker pull $(ALPINE)
	docker pull $(KIND)
	docker pull $(POSTGRES)
	docker pull $(VAULT)
	docker pull $(GRAFANA)
	docker pull $(PROMETHEUS)
	docker pull $(TEMPO)
	docker pull $(TELEPRESENCE)
	
# ======================================================================
# Run Locally

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

# ======================================================================
# Building containers

VERSION := 1.0 

all: sales-api

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ======================================================================
# Running from within k8s/kind

KIND_CLUSTER := jnk-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=sales-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-sales:
	kubectl get pods -o wide --watch --namespace=sales-system

kind-load:
	cd zarf/k8s/kind/sales-pod; kustomize edit set images sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -

kind-delete:
	cat zarf/k8s/base/sales-pod/base-sales.yaml | kubectl delete -f -

kind-restart:
	kubectl rollout restart deployment sales-pod --namespace=sales-system

kind-logs:
	kubectl logs -l app=sales --namespace=sales-system --all-containers -f --tail=100 | go run app/tooling/logfmt/main.go

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=sales --namespace=sales-system

tidy:
	go mod tidy
	go mod vendor