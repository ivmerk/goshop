module github.com/ivmerk/goshop/inventory

go 1.27.3

replace github.com/ivmerk/goshop/shared => ../shared

require (
	github.com/klauspost/compress v1.20.1 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
)

require (
	github.com/caarlos0/env/v11 v11.4.1
	github.com/google/uuid v1.6.0
	github.com/ivmerk/goshop/shared v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.12.1
	go.mongodb.org/mongo-driver/v2 v2.9.2
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.12
)

replace github.com/ivmerk/goshop/platform => ../platform
