GEN=gen
VND=vnd
GENGOPATH=$(shell pwd)/$(GEN)
VNDGOPATH=$(shell pwd)/$(VND)
GOPATH=$(VNDGOPATH):$(shell pwd):$(GENGOPATH)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")
PID=/tmp/go-$(GONAME).pid

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
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get .

install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

run:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFILES)

watch:
	@$(MAKE) restart &
	@fswatch -o . -e 'bin/.*' | xargs -n1 -I{}  make restart

restart: clear stop clean build start

start:
	@echo "Starting bin/$(GONAME)"
	@./bin/$(GONAME) & echo $$! > $(PID)

stop:
	@echo "Stopping bin/$(GONAME) if it's running"
	@-kill `[[ -f $(PID) ]] && cat $(PID)` 2>/dev/null || true

clear:
	@clear

clean:
	@echo "Cleaning"
	@rm -rf $(GENGOPATH)
	@rm -rf $(VNDGOPATH)
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY: build get install run watch start stop restart clean protoc
