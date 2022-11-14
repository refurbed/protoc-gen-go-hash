install-proto-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

build:
	go build ./cmd/protoc-gen-go-hash

test: build
	go test ./...

example: build
	protoc --plugin=protoc-gen-go-hash=./protoc-gen-go-hash --go_out=. --go_opt=paths=source_relative --go-hash_out=. --go-hash_opt=paths=source_relative example/*.proto

clean:
	rm example/example.pb.go example/example_hash.pb.go protoc-gen-go-hash
