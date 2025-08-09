package grpc

import (
	"context"

	pb "github.com/Vin-Xi/auth/gen/token"
	"github.com/Vin-Xi/auth/internal/service"
	"github.com/Vin-Xi/auth/internal/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenServer struct {
	pb.UnimplementedAuthServiceServer
	userService service.Service
	jwtEngine   *util.JWTEngine
}

func NewTokenServer(service service.Service, jwtEngine *util.JWTEngine) *TokenServer {
	return &TokenServer{
		userService: service,
		jwtEngine:   jwtEngine,
	}
}

func (e *TokenServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	token := req.GetToken()

	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is required")
	}

	userID, err := e.jwtEngine.Verify(token)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "token is invalid")
	}

	u, err := e.userService.GetUserByID(ctx, userID)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid user")
	}

	if !u.IsActive {
		return nil, status.Error(codes.PermissionDenied, "user is not active")
	}

	response := &pb.VerifyTokenResponse{
		UserId:   u.ID.String(),
		Email:    u.Email,
		IsActive: u.IsActive,
	}

	return response, nil
}
