all: gen-proto test

test:
	go test ./...

benchmark:
	go test ./... -bench . --benchmem

.PHONY: gen-proto
gen-proto: internal/examples/simple/gogo/gen/logs/logs.pb.go internal/examples/simple/google/gen/logs/logs.pb.go

internal/examples/simple/gogo/gen/logs/logs.pb.go: internal/examples/simple/logs.proto Makefile
	#cd internal && protoc examples/simple/logs.proto --go_out=.
	docker run --rm -v${PWD}:${PWD} \
            -w${PWD} otel/build-protobuf:latest --proto_path=${PWD}/internal/examples/simple \
            --gogofaster_out=plugins=grpc:./internal/examples/simple/gogo/ ${PWD}/internal/examples/simple/logs.proto

internal/examples/simple/google/gen/logs/logs.pb.go: internal/examples/simple/logs.proto Makefile
	docker run --rm -v${PWD}:${PWD} \
            -w${PWD} otel/build-protobuf:latest --proto_path=${PWD}/internal/examples/simple \
            --go_out=plugins=grpc:./internal/examples/simple/google/ ${PWD}/internal/examples/simple/logs.proto
