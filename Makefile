RUNTIME ?= podman
REPOOWNER ?= fromani
IMAGENAME_NUMACELL ?= k8s-dp-numacell
IMAGETAG ?= latest

BUILDFLAGS=GO111MODULE=on GOPROXY=off GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 CGO_ENABLED=0

all: plugins

outdir:
	mkdir -p _output || :

clean:
	rm -rf _output

.PHONY: push
push: push-numacell


# alias
.PHONY: dist
dist: plugins

.PHONY: images
images: image-numacell

.PHONY: push
push: push-numacell

.PHONY: plugins
plugins: build-numacell

build-numacell: outdir
	$(BUILDFLAGS) go build -v -o _output/numacell ./cmd/numacell

.PHONY: image-numacell
image-numacell: build-numacell
	@echo "building numacell image"
	$(RUNTIME) build -f images/Dockerfile.numacell -t quay.io/$(REPOOWNER)/$(IMAGENAME_NUMACELL):$(IMAGETAG) .

.PHONY: push-numacell
push-numacell: image-numacell
	@echo "pushing numacell image"
	$(RUNTIME) push quay.io/$(REPOOWNER)/$(IMAGENAME_NUMACELL):$(IMAGETAG)

