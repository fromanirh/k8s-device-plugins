RUNTIME ?= podman
REPOOWNER ?= fromani
IMAGENAME ?= k8s-device-plugins
IMAGETAG ?= latest

BUILDFLAGS=GO111MODULE=on GOPROXY=off GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 CGO_ENABLED=0

all: dist

outdir:
	mkdir -p _output || :

.PHONY: dist
dist: plugins

.PHONY: plugins
plugins: numacell

numacell: outdir
	$(BUILDFLAGS) go build -v -o _output/numacell ./cmd/numacell

clean:
	rm -rf _output

.PHONY: image
image: plugins
	@echo "building image"
	$(RUNTIME) build -f Dockerfile -t quay.io/$(REPOOWNER)/$(IMAGENAME):$(IMAGETAG) .

.PHONY: push
push: image
	@echo "pushing image"
	$(RUNTIME) push quay.io/$(REPOOWNER)/$(IMAGENAME):$(IMAGETAG)

