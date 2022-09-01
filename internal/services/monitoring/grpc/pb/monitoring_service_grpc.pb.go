// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: monitoring_service.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MonitoringServiceClient is the client API for MonitoringService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MonitoringServiceClient interface {
	Report(ctx context.Context, in *TelemetryRequest, opts ...grpc.CallOption) (*TelemetryResponse, error)
}

type monitoringServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMonitoringServiceClient(cc grpc.ClientConnInterface) MonitoringServiceClient {
	return &monitoringServiceClient{cc}
}

func (c *monitoringServiceClient) Report(ctx context.Context, in *TelemetryRequest, opts ...grpc.CallOption) (*TelemetryResponse, error) {
	out := new(TelemetryResponse)
	err := c.cc.Invoke(ctx, "/monitoring.MonitoringService/Report", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MonitoringServiceServer is the server API for MonitoringService service.
// All implementations must embed UnimplementedMonitoringServiceServer
// for forward compatibility
type MonitoringServiceServer interface {
	Report(context.Context, *TelemetryRequest) (*TelemetryResponse, error)
	mustEmbedUnimplementedMonitoringServiceServer()
}

// UnimplementedMonitoringServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMonitoringServiceServer struct {
}

func (UnimplementedMonitoringServiceServer) Report(context.Context, *TelemetryRequest) (*TelemetryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Report not implemented")
}
func (UnimplementedMonitoringServiceServer) mustEmbedUnimplementedMonitoringServiceServer() {}

// UnsafeMonitoringServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MonitoringServiceServer will
// result in compilation errors.
type UnsafeMonitoringServiceServer interface {
	mustEmbedUnimplementedMonitoringServiceServer()
}

func RegisterMonitoringServiceServer(s grpc.ServiceRegistrar, srv MonitoringServiceServer) {
	s.RegisterService(&MonitoringService_ServiceDesc, srv)
}

func _MonitoringService_Report_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TelemetryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitoringServiceServer).Report(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/monitoring.MonitoringService/Report",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitoringServiceServer).Report(ctx, req.(*TelemetryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MonitoringService_ServiceDesc is the grpc.ServiceDesc for MonitoringService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MonitoringService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "monitoring.MonitoringService",
	HandlerType: (*MonitoringServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Report",
			Handler:    _MonitoringService_Report_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "monitoring_service.proto",
}