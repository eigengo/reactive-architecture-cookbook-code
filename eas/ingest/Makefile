GEN=gen
VND=vnd
GENGOPATH=$(shell pwd)/$(GEN)
VNDGOPATH=$(shell pwd)/$(VND)
GOPATH=$(VNDGOPATH):$(shell pwd):$(GENGOPATH)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")
PID=/tmp/go-$(GONAME).pid

PACKAGES=$(shell find src -type f -name "*.go" -exec dirname {} \; | sed 's|src/||' | sort | uniq)

PROTODIR=../protocol
PROTOFILES=$(shell find $(PROTODIR) -type f -name '*.proto')

all: protoc get build

build:
	@echo "Building $(GOFILES) using $(GOPATH) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

protoc:
	@echo "Building protocols from $(PROTOFILES)"
	@GOPATH=$(GOPATH) go get -u github.com/golang/protobuf/protoc-gen-go
	mkdir -p $(GENGOPATH)/src
	$(foreach p,$(PROTOFILES),PATH=$(PATH):$(VNDGOPATH)/bin protoc --go_out=$(GEN)/src $(p) -I $(PROTODIR);)

get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get -t .
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get -t github.com/reactivesystemsarchitecture/eas/ingest

test:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test github.com/reactivesystemsarchitecture/eas/ingest -v -covermode=atomic -o bin/test -coverprofile=bin/cover.out
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go tool cover --html bin/cover.out


clean:
	@echo "Cleaning"
	@rm -rf $(GENGOPATH)
	@rm -rf $(VNDGOPATH)
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY: build get clean protoc
