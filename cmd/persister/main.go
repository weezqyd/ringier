package main

import (
	"context"
	"github.com/weezqyd/ringier/internal/adapters/database"
	gRPC "github.com/weezqyd/ringier/internal/services/persist/grpc"
	"log"
)

func main() {
	ctx := context.Background()
	db, err := database.NewMongoClient(
		ctx,
		"mongodb://root:secret@mongodb:27017/",
		"ringier_events",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDbConnection(ctx)
	log.Println("[persister] Starting gRPC server on port 9090")
	server := gRPC.NewService(db)
	server.Run(":9090")
}
