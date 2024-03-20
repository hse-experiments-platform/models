package models

import (
	"context"

	pb "github.com/hse-experiments-platform/models/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *modelsService) GetFullTrainedModel(context.Context, *pb.GetFullTrainedModelRequest) (*pb.GetFullTrainedModelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFullTrainedModel not implemented")
}
