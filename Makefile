all: clean tools generate test build gosec

clean:
	rm -rf build/
	rm -rf internal/rpc/payment/*.pb.go

generate:
	go generate proto/generate.go
	find . -type f -name 'generate_test.go' -exec \
	go generate {} \;

test:
	go test ./... -count=1 -coverprofile test-coverage.out

build:
	go build -o build/ab-bsn ./cmd/bsn/

gosec:
	gosec -fmt=sonarqube -out gosec_report.json -no-fail ./...

tools:
	go install github.com/vektra/mockery/v2
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: \
	all
	clean
	generate
	tools
	test
	build
	gosec
