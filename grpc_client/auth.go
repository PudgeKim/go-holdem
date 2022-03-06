package grpc_client

import (
	"context"
	pb "github.com/PudgeKim/go-holdem-protos/protos"
)

func (h *GrpcHandler) GetUser(ctx context.Context, userId string) (*pb.User, error) {
	user, err := h.client.GetUser(ctx, &pb.UserId{Id: userId})
	if err != nil {
		return nil, err
	}

	return user, nil
}
