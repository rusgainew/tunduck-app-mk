module github.com/rusgainew/tunduck-app-mk/auth-service

go 1.25

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/lib/pq v1.10.9
	github.com/rabbitmq/amqp091-go v1.9.0
	github.com/redis/go-redis/v9 v9.17.2
	github.com/rusgainew/tunduck-app-mk/proto-lib v0.0.0
	golang.org/x/crypto v0.46.0
	google.golang.org/grpc v1.68.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	go.uber.org/goleak v1.3.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/rusgainew/tunduck-app-mk/proto-lib => ../proto-lib

exclude google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1
