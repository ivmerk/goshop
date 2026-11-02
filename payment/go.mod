module github.com/ivmerk/goshop/payment

go 1.27.3

require (
	github.com/stretchr/objx v0.5.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/ivmerk/goshop/shared => ../shared

require (
	github.com/caarlos0/env/v11 v11.4.1
	github.com/google/uuid v1.6.0
	github.com/ivmerk/goshop/platform v0.0.0
	github.com/ivmerk/goshop/shared v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.12.1
	go.uber.org/zap v1.28.0
	google.golang.org/grpc v1.84.0
)

replace github.com/ivmerk/goshop/platform => ../platform
