export TARGET_HOME=home
export TARGET_ARTICLE=article
export TARGET_USERRPC=userrpc
export TARGET_USER=user
export DOCKER_TARGET=hub.xxx.com/community-blogger
export DOCKER_BUILDER_TARGET=$(DOCKER_TARGET).builder
apps = 'home' 'article' 'userrpc' 'user'

.PHONY: build
DOCKER_TAG := $(if $(DOCKER_TAG),$(DOCKER_TAG),latest)
build:
	go build -o ./$(TARGET_HOME) ./cmd/$(TARGET_HOME)/
	go build -o ./$(TARGET_ARTICLE) ./cmd/$(TARGET_ARTICLE)/
	go build -o ./$(TARGET_USERRPC) ./cmd/$(TARGET_USERRPC)/
	go build -o ./$(TARGET_USER) ./cmd/$(TARGET_USER)/
.PHONY: docker-build
docker-build:
	DOCKER_BUILDKIT=0 docker build --build-arg TARGET_HOME --build-arg TARGET_ARTICLE --build-arg TARGET_USERRPC --build-arg TARGET_USER --build-arg GOPRYXY	--target builder -t $(DOCKER_BUILDER_TARGET) -f Dockerfile .
.PHONY: docker
docker:
	DOCKER_BUILDKIT=0 docker build --build-arg TARGET_HOME --build-arg TARGET_ARTICLE --build-arg TARGET_USERRPC --build-arg TARGET_USER -t $(DOCKER_TARGET) -f Dockerfile .
	DOCKER_BUILDKIT=0 docker build --build-arg TARGET_HOME --build-arg TARGET_ARTICLE --build-arg TARGET_USERRPC --build-arg TARGET_USER -t $(DOCKER_TARGET):dev -f Dockerfile .
docker-release:
	docker push $(DOCKER_BUILDER_TARGET)
	docker push $(DOCKER_TARGET):$(DOCKER_TAG)
docker-push:
	docker push $(DOCKER_BUILDER_TARGET)
	docker push $(DOCKER_TARGET):dev

.PHONY: wire
wire:
	wire ./...