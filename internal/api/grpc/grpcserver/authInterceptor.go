package grpcserver

import (
	"context"
	"encoding/base64"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextValue string

func (s Server) BasicAuthInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	isAdmin := false

	md, exist := metadata.FromIncomingContext(ctx)
	switch exist {
	case true:
		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			break
		}

		header, found := strings.CutPrefix(authHeader[0], "Basic ")
		if !found {
			break
		}

		credentials, err := base64.StdEncoding.DecodeString(header)
		if err != nil {
			break
		}

		parts := strings.SplitN(string(credentials), ":", 2)
		if len(parts) != 2 {
			break
		}

		username := parts[0]
		password := parts[1]

		result, err := s.service.CheckPassword(ctx, username, password)
		if err != nil {
			break
		} else if !result {
			break
		}

		isAdmin = true
	default:
	}

	ctx = context.WithValue(ctx, contextValue("isAdmin"), isAdmin)

	return handler(ctx, req)
}
