SHELL := /bin/bash

# ============================================================================
# Testing running system ( run this on cmd line while app running )

# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

# ============================================================================

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

<<<<<<< HEAD
# ======================================================================
# Building containers
=======
>>>>>>> set up debugging support for use during development

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