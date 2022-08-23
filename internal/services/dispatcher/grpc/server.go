package grpc

import (
	"context"
	"fmt"
	common "github.com/weezqyd/ringier/internal/common/messages"
	"github.com/weezqyd/ringier/internal/services/dispatcher/grpc/dispatcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"net"
)

type Service struct {
	dispatcher.UnimplementedDispatcherServiceServer
}

// Run registers the Dispatcher service to a grpcServer and serves on
// the specified port
func (s *Service) Run(port string) {
	var err error

	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s %v", port, err)
	}
	sever := s
	grpcServer := grpc.NewServer()
	dispatcher.RegisterDispatcherServiceServer(grpcServer, sever)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve gRPC server over port %s %v", port, err)
	}
}

func (s *Service) Dispatch(ctx context.Context, req *common.EventRequest) (*common.EventResponse, error) {
	resp := &common.EventResponse{
		RequestId: req.RequestId,
		Errors:    []string{},
	}
	// nothing fancy here just return an error if the random value is less than 20%
	if err := rand.Intn(100); err <= int(20) {
		resp.Success = false
		resp.Errors = []string{fmt.Sprintf("failed to dispatch the event")}

		return resp, status.Error(codes.Internal, "internal error")
	}
	resp.Success = true

	return resp, nil
}
