module github.com/hamza4253/tiny-url/gateway

go 1.23.0

replace github.com/hamza4253/tiny-url/shared => ../shared

require (
	github.com/hamza4253/tiny-url/shared v0.0.0-00010101000000-000000000000
	github.com/matoous/go-nanoid/v2 v2.1.0
	github.com/streadway/amqp v1.1.0
	google.golang.org/grpc v1.75.0
)

require (
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.7 // indirect
)
