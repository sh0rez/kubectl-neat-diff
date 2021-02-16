.PHONY: lint test static install uninstall cross
BIN_DIR := $(GOPATH)/bin
GOX := $(BIN_DIR)/gox

LDFLAGS := '-s -w -extldflags "-static"'
static:
	CGO_ENABLED=0 go build -ldflags=${LDFLAGS} .

install:
	CGO_ENABLED=0 go install -ldflags=${LDFLAGS} .

$(GOX):
	go get -u github.com/mitchellh/gox
cross: $(GOX)
	CGO_ENABLED=0 gox -output="dist/{{.Dir}}-{{.OS}}-{{.Arch}}" -ldflags=${LDFLAGS} -osarch="linux/amd64" -osarch="darwin/amd64" -osarch="windows/amd64" .
