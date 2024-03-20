package models

import (
	"context"

	pb "github.com/hse-experiments-platform/models/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *modelsService) GetTrainedModels(context.Context, *pb.GetTrainedModelsRequest) (*pb.GetTrainedModelsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTrainedModels not implemented")
}
