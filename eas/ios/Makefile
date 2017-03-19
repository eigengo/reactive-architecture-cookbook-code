GEN=EAS\ Protocol
GENPATH=$(shell pwd)/$(GEN)

PROTODIR=../protocol
PROTOFILES=$(shell find $(PROTODIR) -type f -name '*.proto')

all: protoc

protoc:
	# Use https://github.com/alexeyxo/protobuf-swift
	@echo "Building protocols from $(PROTOFILES)"
	mkdir -p $(GENPATH)
	$(foreach p,$(PROTOFILES),PATH=$(PATH) protoc --swift_out=$(GEN) $(p) -I $(PROTODIR);)

.PHONY: protoc
