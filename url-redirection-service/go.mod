module github.com/hamza4253/tiny-url/redirect

go 1.23.0

require (
	github.com/hamza4253/tiny-url/shared v0.0.0-00010101000000-000000000000
	github.com/redis/go-redis/v9 v9.12.1
	google.golang.org/grpc v1.74.2
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.7 // indirect
)

replace github.com/hamza4253/tiny-url/shared => ../shared
