package main

import (
	gRPC "github.com/weezqyd/ringier/internal/services/dispatcher/grpc"
	"log"
)

func main() {
	log.Println("[dispatcher] Starting gRPC server on port 9090")
	server := gRPC.Service{}
	server.Run(":9090")
}
