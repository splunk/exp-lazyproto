all: gen-proto test

test:
	go test ./...

benchmark:
	go test ./... -bench .

.PHONY: gen-proto
gen-proto: internal/examples/simple/google.logs.pb.go

internal/examples/simple/google.logs.pb.go: internal/examples/simple/logs.proto
	cd internal && protoc examples/simple/logs.proto --go_out=.