package grpc_client

import (
	pb "github.com/PudgeKim/go-holdem-protos/protos"
)

type GrpcHandler struct {
	client pb.AuthClient
}

func NewGrpcHandler(client pb.AuthClient) *GrpcHandler {
	return &GrpcHandler{
		client: client,
	}
}
