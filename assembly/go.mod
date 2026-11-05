module github.com/ivmerk/goshop/assembly

go 1.27.3

require (
	github.com/IBM/sarama v1.61.1
	github.com/caarlos0/env/v11 v11.4.2
	github.com/google/uuid v1.6.0
	github.com/ivmerk/goshop/payment v0.0.0-20261104121609-a352321f3274
	github.com/ivmerk/goshop/platform v0.0.0
	github.com/ivmerk/goshop/shared v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.12.1
	go.uber.org/zap v1.28.0
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.20.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.33 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20250401214520-65e299d6c5c9 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
	google.golang.org/grpc v1.84.0 // indirect
)

replace github.com/ivmerk/goshop/shared => ../shared

replace github.com/ivmerk/goshop/platform => ../platform
