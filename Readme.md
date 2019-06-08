# gRPC poc

Testing grpc client and server


## Running

### Generate proto

```cmd
 make gen-proto
```

### Generate keys

```cmd
 make gen-server-tls
```

### Run server

```cmd
 make run-server
```

### Run client

```cmd
 make run-client
```

## Requirements

### gRPC

```cmd
go get -u google.golang.org/grpc
```

### Protocol Buffers v3

```cmd
go get -u github.com/golang/protobuf/protoc-gen-go
```

### Protoc plugin for go

```cmd
go get -u github.com/golang/protobuf/protoc-gen-go
```

### Export gobin

```cmd
export PATH=$PATH:$GOPATH/bin
```
