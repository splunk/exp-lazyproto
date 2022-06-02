module github.com/tigrannajaryan/exp-lazyproto

go 1.18

require (
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.0
	github.com/jhump/protoreflect v1.12.0
	github.com/stretchr/testify v1.7.1
	github.com/tigrannajaryan/molecule v1.0.0
	github.com/yoheimuta/go-protoparser/v4 v4.6.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace github.com/tigrannajaryan/molecule => ../molecule
