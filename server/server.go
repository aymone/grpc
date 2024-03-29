package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/aymone/grpc/handler"
	pb "github.com/aymone/grpc/proto/service"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Run server
func Run(addr, clientAddr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to initializa TCP listen: %v", err)
	}

	defer lis.Close()
	go runGRPC(lis, handler.New())
	runHTTP(clientAddr)
}

func runGRPC(lis net.Listener, handler pb.SimpleServerServer) {
	creds, err := credentials.NewServerTLSFromFile("server/server-cert.pem", "server/server-key.pem")
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}

	server := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(recoveryInterceptor),
	)

	pb.RegisterSimpleServerServer(server, handler)

	log.Printf("gRPC Listening on %s\n", lis.Addr().String())

	server.Serve(lis)
}

func runHTTP(clientAddr string) {
	runtime.HTTPError = customHTTPError

	addr := ":6001"
	creds, err := credentials.NewClientTLSFromFile("server/server-cert.pem", "")
	if err != nil {
		log.Fatalf("gateway cert load error: %s", err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	mux := runtime.NewServeMux()
	if err := pb.RegisterSimpleServerHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
	log.Printf("HTTP Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func recoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from", r)
		}
	}()
	return handler(ctx, req)
}

// func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	meta, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
// 	}
// 	if len(meta["authorization"]) != 1 {
// 		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
// 	}
// 	if meta["authorization"][0] != "valid-token" {
// 		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
// 	}

// 	return handler(ctx, req)
// }

type errorBody struct {
	Err string `json:"error,omitempty"`
}

func customHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))
	jErr := json.NewEncoder(w).Encode(errorBody{
		Err: grpc.ErrorDesc(err),
	})

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}
