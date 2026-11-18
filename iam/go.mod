module github.com/ivmerk/goshop/iam

go 1.27.3

require (
	github.com/caarlos0/env/v11 v11.4.2
	github.com/gomodule/redigo v1.9.3
	github.com/google/uuid v1.6.0
	github.com/ivmerk/goshop/platform v0.0.0-00010101000000-000000000000
	github.com/ivmerk/goshop/shared v0.0.0
	github.com/jackc/pgx/v5 v5.13.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.58.0
	google.golang.org/grpc v1.84.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/ivmerk/goshop/shared => ../shared

replace github.com/ivmerk/goshop/platform => ../platform
