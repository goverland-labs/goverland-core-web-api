module github.com/goverland-labs/goverland-core-web-api

go 1.23

toolchain go1.23.1

replace github.com/goverland-labs/goverland-core-web-api/protocol => ./protocol

require (
	github.com/caarlos0/env/v10 v10.0.0
	github.com/golang/protobuf v1.5.4
	github.com/google/uuid v1.6.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/goverland-labs/goverland-core-feed/protocol v0.2.1
	github.com/goverland-labs/goverland-core-storage/protocol v0.4.20-0.20250624151607-3a5e299a5521
	github.com/goverland-labs/goverland-core-web-api/protocol v0.0.0-20250220134513-ce50ab1484b8
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/prometheus/client_golang v1.14.0
	github.com/rs/zerolog v1.29.0
	github.com/s-larionov/process-manager v0.0.1
	github.com/shopspring/decimal v1.3.1
	go.openly.dev/pointy v1.3.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.3
)

require (
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
