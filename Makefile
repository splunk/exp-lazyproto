internal/examples/simple/google.logs.pb.go: internal/examples/simple/logs.proto
	cd internal && protoc examples/simple/logs.proto --go_out=.