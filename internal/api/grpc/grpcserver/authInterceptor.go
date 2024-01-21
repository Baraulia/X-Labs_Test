package grpcserver

import (
	"context"
	"encoding/base64"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//nolint:lll
func (s Server) BasicAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	authorizedMethods := map[string]struct{}{
		"/user.UserService/CreateUser": {},
		"/user.UserService/UpdateUser": {},
		"/user.UserService/DeleteUser": {},
	}

	_, exist := authorizedMethods[info.FullMethod]
	switch exist {
	case true:
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			s.logger.Error("BasicAuthInterceptor: missing metadata", nil)
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			s.logger.Error("BasicAuthInterceptor: missing Authorization header", nil)
			return nil, status.Error(codes.Unauthenticated, "missing Authorization header")
		}

		header, ok := strings.CutPrefix(authHeader[0], "Basic ")
		if !ok {
			s.logger.Error("BasicAuthInterceptor: invalid Authorization header", nil)
			return nil, status.Error(codes.Unauthenticated, "invalid Authorization header")
		}

		credentials, err := base64.StdEncoding.DecodeString(header)
		if err != nil {
			s.logger.Error("BasicAuthInterceptor: invalid base64 encoding", nil)
			return nil, status.Error(codes.Unauthenticated, "invalid base64 encoding")
		}

		parts := strings.SplitN(string(credentials), ":", 2)
		if len(parts) != 2 {
			s.logger.Error("BasicAuthInterceptor: invalid credentials format", nil)
			return nil, status.Error(codes.Unauthenticated, "invalid credentials format")
		}

		username := parts[0]
		password := parts[1]

		result, err := s.service.CheckPassword(ctx, username, password)
		if err != nil {
			s.logger.Error("BasicAuthInterceptor: error while checking password", map[string]interface{}{"error": err})
			return nil, status.Error(codes.Unauthenticated, err.Error())
		} else if !result {
			s.logger.Error("BasicAuthInterceptor: password mismatch", nil)
			return nil, status.Error(codes.Unauthenticated, "password mismatch")
		}
	default:
	}

	return handler(ctx, req)
}
