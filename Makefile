HELM_PLUGIN_DIR ?= $(shell helm env | grep HELM_PLUGINS | cut -d\" -f2)/helm-ssm
HELM_PLUGIN_NAME := helm-ssm
VERSION := $(shell cat .version)
DIST := $(CURDIR)/_dist
LDFLAGS := "-X main.version=${VERSION}"

.PHONY: install
install: dist
	@if [ ! -f .version ] ; then echo "dev" > .version ; fi
	mkdir -p $(HELM_PLUGIN_DIR)
	@if [ "$$(uname)" = "Darwin" ]; then file="${HELM_PLUGIN_NAME}-macos"; \
 	elif [ "$$(uname)" = "Linux" ]; then file="${HELM_PLUGIN_NAME}-linux"; \
	else file="${HELM_PLUGIN_NAME}-windows"; \
	fi; \
	mkdir -p $(DIST)/$$file ; \
	tar -xf $(DIST)/$$file.tgz -C $(DIST)/$$file ; \
	cp -r $(DIST)/$$file/* $(HELM_PLUGIN_DIR) ;\
	rm -rf $(DIST)/$$file

.PHONY: hookInstall
hookInstall: build

.PHONY: build
build:
	go build -o bin/${HELM_PLUGIN_NAME} -ldflags $(LDFLAGS) ./cmd

.PHONY: test
test:
	go test -v ./internal

.PHONY: dist
dist:
	mkdir -p $(DIST)
	sed -i.bak 's/version:.*/version: "'$(VERSION)'"/g' plugin.yaml && rm plugin.yaml.bak
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${HELM_PLUGIN_NAME} -ldflags $(LDFLAGS) ./cmd
	tar -zcvf $(DIST)/${HELM_PLUGIN_NAME}-linux.tgz ${HELM_PLUGIN_NAME} README.md LICENSE plugin.yaml
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${HELM_PLUGIN_NAME} -ldflags $(LDFLAGS) ./cmd
	tar -zcvf $(DIST)/${HELM_PLUGIN_NAME}-macos.tgz ${HELM_PLUGIN_NAME} README.md LICENSE plugin.yaml
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${HELM_PLUGIN_NAME}.exe -ldflags $(LDFLAGS) ./cmd
	tar -zcvf $(DIST)/${HELM_PLUGIN_NAME}-windows.tgz ${HELM_PLUGIN_NAME}.exe README.md LICENSE plugin.yaml
	rm ${HELM_PLUGIN_NAME}
	rm ${HELM_PLUGIN_NAME}.exe