module github.com/linemk/rocket-shop/payment

go 1.24.7

require (
	github.com/google/uuid v1.6.0
	github.com/linemk/rocket-shop/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.68.0
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/protobuf v1.36.0 // indirect
)

replace github.com/linemk/rocket-shop/shared => ../shared
