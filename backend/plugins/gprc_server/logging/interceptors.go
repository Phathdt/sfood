package logging

import (
	"context"
	"time"

	sctx "github.com/viettranx/service-context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(log sctx.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		startTime := time.Now()
		result, err := handler(ctx, req)
		duration := time.Since(startTime)

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}
		log.Withs(sctx.Fields{
			"protocol":    "grpc",
			"method":      info.FullMethod,
			"status_code": int(statusCode),
			"status_text": statusCode.String(),
			"duration":    duration,
		}).Info("received a gRPC request")

		return result, err
	}
}

func StreamServerInterceptor(log sctx.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		startTime := time.Now()
		err = handler(srv, stream)
		duration := time.Since(startTime)

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		log.Withs(sctx.Fields{
			"protocol":    "grpc",
			"method":      info.FullMethod,
			"status_code": int(statusCode),
			"status_text": statusCode.String(),
			"duration":    duration,
		}).Info("received a gRPC request")

		return err
	}
}
