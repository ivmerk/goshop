module github.com/ivmerk/goshop/payment

go 1.27.3

require (
	github.com/stretchr/objx v0.5.3 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/ivmerk/goshop/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/ivmerk/goshop/shared v0.0.0
	github.com/stretchr/testify v1.12.1
	google.golang.org/grpc v1.84.0
)

replace github.com/ivmerk/goshop/platform => ../platform
