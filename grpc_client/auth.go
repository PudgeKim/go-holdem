package grpc_client

import (
	"context"
	pb "github.com/PudgeKim/go-holdem-protos/protos"
	"github.com/PudgeKim/go-holdem/grpc_client/grpc_error"
)

func (h *GrpcHandler) GetUser(ctx context.Context, userId string) (*pb.User, error) {
	user, err := h.client.GetUser(ctx, &pb.UserId{Id: userId})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, grpc_error.UserNotExist
	}

	return user, nil
}
