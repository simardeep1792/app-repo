.PHONY: build push toggle-bad toggle-good bump-dev

IMAGE_REGISTRY ?= ghcr.io
IMAGE_REPO ?= $(shell echo $${GITHUB_REPOSITORY:-simardeep1792/app-repo})
IMAGE_TAG ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "latest")
IMAGE := $(IMAGE_REGISTRY)/$(IMAGE_REPO):$(IMAGE_TAG)

build:
	docker build -t $(IMAGE) .
	@echo "Built image: $(IMAGE)"

push: build
	docker push $(IMAGE)
	@echo "Pushed image: $(IMAGE)"

toggle-bad:
	@echo "Enabling failure injection"
	@sed -i.bak 's/INJECT_FAILURE=false/INJECT_FAILURE=true/' Dockerfile || true
	@grep -q "ENV INJECT_FAILURE" Dockerfile || echo "ENV INJECT_FAILURE=true" >> Dockerfile
	@echo "Failure injection enabled. Commit and push to trigger bad release."

toggle-good:
	@echo "Disabling failure injection"
	@sed -i.bak '/ENV INJECT_FAILURE/d' Dockerfile || true
	@rm -f Dockerfile.bak
	@echo "Failure injection disabled. Commit and push to trigger good release."

bump-dev:
	@echo "Updating dev environment with image: $(IMAGE)"
	@cd ../env-repo && \
		yq e '.spec.template.spec.containers[0].image = "$(IMAGE)"' -i envs/dev/rollout.yaml && \
		echo "Updated envs/dev/rollout.yaml with image $(IMAGE)"