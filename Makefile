HELM_HOME ?= $(shell helm home)
HELM_PLUGIN_DIR ?= $(HELM_HOME)/plugins/helm-ssm
HELM_PLUGIN_NAME := helm-ssm
HAS_DEP := $(shell which dep)
DEP_VERSION := v0.5.1
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist
LDFLAGS := "-X main.version=${VERSION}"

.PHONY: install
install: bootstrap dist
	mkdir -p $(HELM_PLUGIN_DIR)
	@if [ "$$(uname)" = "Darwin" ]; then file="${HELM_PLUGIN_NAME}-macos-${VERSION}"; \
 	elif [ "$$(uname)" = "Linux" ]; then file="${HELM_PLUGIN_NAME}-linux-${VERSION}"; \
	else file="${HELM_PLUGIN_NAME}-windows-${VERSION}"; \
	fi; \
	mkdir -p $(DIST)/$$file ; \
	tar -xf $(DIST)/$$file.tgz -C $(DIST)/$$file ; \
	cp -r $(DIST)/$$file/* $(HELM_PLUGIN_DIR) ;\
	rm -rf $(DIST)/$$file

.PHONY: hookInstall
hookInstall: bootstrap build

.PHONY: build
build:
	go build -o bin/${HELM_PLUGIN_NAME} -ldflags $(LDFLAGS) ./cmd

.PHONY: test
test: bootstrap
	go test -v ./internal

.PHONY: dist
dist:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${HELM_PLUGIN_NAME} -ldflags $(LDFLAGS) ./cmd
	chmod +x ${HELM_PLUGIN_NAME}
	tar -zcvf $(DIST)/${HELM_PLUGIN_NAME}-linux-$(VERSION).tgz ${HELM_PLUGIN_NAME} README.md LICENSE plugin.yaml
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${HELM_PLUGIN_NAME} -ldflags $(LDFLAGS) ./cmd
	chmod +x ${HELM_PLUGIN_NAME}
	tar -zcvf $(DIST)/${HELM_PLUGIN_NAME}-macos-$(VERSION).tgz ${HELM_PLUGIN_NAME} README.md LICENSE plugin.yaml
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${HELM_PLUGIN_NAME}.exe -ldflags $(LDFLAGS) ./cmd
	chmod +x ${HELM_PLUGIN_NAME}
	tar -zcvf $(DIST)/${HELM_PLUGIN_NAME}-windows-$(VERSION).tgz ${HELM_PLUGIN_NAME}.exe README.md LICENSE plugin.yaml
	rm ${HELM_PLUGIN_NAME}
	rm ${HELM_PLUGIN_NAME}.exe

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	DEP_RELEASE_TAG=$(DEP_VERSION) curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	chmod +x $(GOPATH)/bin/dep
endif
	dep ensure