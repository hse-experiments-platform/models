package models

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/hse-experiments-platform/models/pkg/models"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *modelsService) GetFullTrainedModel(ctx context.Context, req *pb.GetFullTrainedModelRequest) (*pb.GetFullTrainedModelResponse, error) {
	resp := &pb.GetFullTrainedModelResponse{Model: &pb.TrainedModel{}}

	modelID := int64(req.GetTrainedModelID())

	model, err := s.commonDB.GetTrainedModel(ctx, modelID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("model with id %v not found", modelID))
	} else if err != nil {
		return nil, fmt.Errorf("txdb.GetModel: %w", err)
	}

	resp.Model = &pb.TrainedModel{
		TrainedModelID: uint64(model.ID),
		Name:           model.Name,
		TrainStatus:    convertTrainingStatus(model.LaunchStatus),
		BaseModelID:    uint64(model.ModelID),
		BaseModelName:  model.ModelName,
		Problem: &pb.ShortProblem{
			Id:          uint64(model.ProblemID),
			Name:        model.ProblemName,
			Description: model.ProblemDescription,
		},
		TrainDatasetID:   uint64(model.TrainDatasetID),
		TrainDatasetName: model.TrainingDatasetName,
		CreatedAt:        timestamppb.New(model.CreatedAt.Time),
		TargetColumn:     model.TargetColumn,
		LaunchID:         uint64(model.LaunchID.Int64),
	}

	return resp, nil
}
