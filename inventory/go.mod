module github.com/linemk/rocket-shop/inventory

go 1.24.7

require (
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/linemk/rocket-shop/platform v0.0.0-00010101000000-000000000000
	github.com/linemk/rocket-shop/shared v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	go.mongodb.org/mongo-driver v1.17.6
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/brianvoe/gofakeit/v7 v7.9.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20250403155104-27863c87afa6 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/onsi/ginkgo/v2 v2.27.2 // indirect
	github.com/onsi/gomega v1.38.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/mod v0.28.0 // indirect
	golang.org/x/net v0.45.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	golang.org/x/tools v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250825161204-c5933d9347a5 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/linemk/rocket-shop/platform => ../platform
	github.com/linemk/rocket-shop/shared => ../shared
)
