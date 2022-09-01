package grpc

import (
	"context"
	"github.com/weezqyd/ringier/internal/services/monitoring/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type Service struct {
	pb.UnimplementedMonitoringServiceServer
}

// Run registers the monitoring service to a grpcServer and serves on
// the specified port
func (s *Service) Run(port string) {
	var err error

	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s %v", port, err)
	}
	sever := s
	grpcServer := grpc.NewServer()
	pb.RegisterMonitoringServiceServer(grpcServer, sever)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve gRPC server over port %s %v", port, err)
	}
}

func (s *Service) Report(context.Context, *pb.TelemetryRequest) (*pb.TelemetryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Report not implemented")
}
