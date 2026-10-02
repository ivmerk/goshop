module github.com/ivmerk/goshop/inventory

go 1.27

replace github.com/ivmerk/goshop/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/ivmerk/goshop/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.12
)

require (
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
)
