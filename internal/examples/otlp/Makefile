GENDIR := gen

# Find all .proto files.
PROTO_FILES := $(wildcard opentelemetry-proto/opentelemetry/proto/*/*/*.proto opentelemetry-proto/opentelemetry/proto/*/*/*/*.proto)

# Function to execute a command. Note the empty line before endef to make sure each command
# gets executed separately instead of concatenated with previous one.
# Accepts command to execute as first parameter.
define exec-command
$(1)

endef

# Generate all implementations
.PHONY: gen-all
gen-all: gen-go

OTEL_DOCKER_PROTOBUF ?= otel/build-protobuf:0.9.0
PROTOC :=   go run ../../../cmd/main.go --proto_path=${PWD} --proto_path=opentelemetry-proto --with_presence

PROTO_GEN_GO_DIR ?= $(GENDIR)/go


# Generate gRPC/Protobuf implementation for Go.
.PHONY: gen-go
gen-go:
	rm -rf ./$(PROTO_GEN_GO_DIR)
	mkdir -p ./$(PROTO_GEN_GO_DIR)
	$(foreach file,$(PROTO_FILES),$(call exec-command,$(PROTOC) --go_out=./$(PROTO_GEN_GO_DIR) $(file)))
