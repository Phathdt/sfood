package gprc_server

import (
	"flag"
	"fmt"
	"net"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	sctx "github.com/viettranx/service-context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"sfood/plugins/gprc_server/logging"
	"sfood/plugins/gprc_server/recovermiddleware"
)

type gprcServer struct {
	id          string
	prefix      string
	port        int
	server      *grpc.Server
	logger      sctx.Logger
	registerHdl func(*grpc.Server)
}

func NewGprcServer(id string) *gprcServer {
	return &gprcServer{id: id}
}

func (s *gprcServer) ID() string {
	return s.id
}

func (s *gprcServer) InitFlags() {
	flag.IntVar(&s.port, "grpc_port", 50051, "Port of gRPC service")
}

func (s *gprcServer) Recover() {
	if err := recover(); err != nil {
		s.logger.Error("recover error", err)
	}
}

func (s *gprcServer) Activate(sc sctx.ServiceContext) error {
	go func() {
		defer s.Recover()

		s.logger = sc.Logger(s.id)

		s.logger.Infoln("Setup gRPC service:", s.prefix)
		s.server = grpc.NewServer(
			grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
				otelgrpc.StreamServerInterceptor(),
				recovermiddleware.StreamServerInterceptor(),
				logging.StreamServerInterceptor(s.logger),
			)),
			grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
				otelgrpc.UnaryServerInterceptor(),
				recovermiddleware.UnaryServerInterceptor(),
				logging.UnaryServerInterceptor(s.logger),
			)),
		)

		if s.registerHdl != nil {
			s.logger.Infoln("registering services...")
			s.registerHdl(s.server)
		}

		address := fmt.Sprintf("0.0.0.0:%d", s.port)
		lis, err := net.Listen("tcp", address)

		if err != nil {
			s.logger.Errorln("Error %v", err)
		}

		_ = s.server.Serve(lis)
	}()

	return nil
}

func (s *gprcServer) Stop() error {
	s.server.Stop()

	return nil
}

func (s *gprcServer) SetRegisterHdl(hdl func(*grpc.Server)) {
	s.registerHdl = hdl
}
