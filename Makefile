IMAGE ?= mkimuram/k8sviz
TAG ?= $(shell cat version.txt)
DEVEL_IMAGE ?= k8sviz
DEVEL_TAG ?= devel

test: test-lint test-fmt test-vet test-unit
	@echo "[Running test]"

test-lint:
	@echo "[Running golint]"
	golint -set_exit_status cmd/... pkg/...

test-fmt:
	@echo "[Running gofmt]"
	if [ "$$(gofmt -l cmd/ pkg/ | wc -l)" -ne 0 ]; then \
		gofmt -d cmd/ pkg/ ;\
		false; \
	fi

test-vet:
	@echo "[Running go vet]"
	go vet `go list ./... | grep -v test/e2e`


test-unit:
	@echo "[Running unit tests]"
	go test -cover `go list ./... | grep -v test/e2e`

test-e2e: build
	@echo "[Running e2e tests]"
	./test/e2e/e2e.sh

build:
	@echo "[Build]"
	mkdir -p bin/
	GO111MODULE=on go build -o bin/k8sviz ./cmd/k8sviz
	mv bin/k8sviz ./

release: test build test-e2e

image-build:
	@echo "[Building image $(DEVEL_IMAGE):$(DEVEL_TAG)]"
	docker build -t $(DEVEL_IMAGE):$(DEVEL_TAG) .

image-push: image-build
	@echo "[Pushing image $(IMAGE):$(TAG)]"
	docker tag $(DEVEL_IMAGE):$(DEVEL_TAG) $(IMAGE):$(TAG)
	docker push $(IMAGE):$(TAG)

.PHONY: test test-lint test-fmt test-vet test-unit test-e2e build release image-build image-push
