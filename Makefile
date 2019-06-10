gen-proto:
	protoc -I. -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:${GOPATH}/src proto/service/service.proto \
		--grpc-gateway_out=logtostderr=true:.

	protoc \
		--go_out=${GOPATH}/src proto/message/message.proto \
		--grpc-gateway_out=logtostderr=true:.

gen-server-tls:
	openssl req -x509 -newkey rsa:4096 -keyout server/server-key.pem -out server/server-cert.pem -days 365 -nodes -subj '/CN=localhost'

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go