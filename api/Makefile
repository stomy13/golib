install:
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go
	go install github.com/bufbuild/buf/cmd/buf
	go install google.golang.org/protobuf/cmd/protoc-gen-go

run:
	go run .

bufcurl:
	buf curl --protocol grpc \
		--http2-prior-knowledge \
    --data '{"number": "23", "text": "hello?"}' \
    http://localhost:8080/connect.ping.v1.PingService/Ping

curl:
	curl \
    --header "Content-Type: application/json" \
    --data '{"number": "23", "text": "hello?"}' \
    http://localhost:8080/connect.ping.v1.PingService/Ping

gen:
	buf generate

lint:
	buf lint

tidy:
	go mod tidy