package main

import (
	"fmt"

	"github.com/aymone/grpc/server"
)

func main() {
	addr := ":6000"
	clientAddr := fmt.Sprintf("localhost%s", addr)
	server.Run(addr, clientAddr)
}
