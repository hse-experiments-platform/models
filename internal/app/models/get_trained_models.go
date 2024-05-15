package models

import (
	"context"
	"fmt"

	"github.com/hse-experiments-platform/models/internal/pkg/storage/db"
	pb "github.com/hse-experiments-platform/models/pkg/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *modelsService) GetTrainedModels(ctx context.Context, req *pb.GetTrainedModelsRequest) (*pb.GetTrainedModelsResponse, error) {
	modelRows, err := s.commonDB.GetTrainedModels(ctx, db.GetTrainedModelsParams{
		Name:    "%" + req.GetQuery() + "%",
		Limit:   int64(req.GetLimit()),
		Offset:  int64(req.GetOffset()),
		ModelID: int64(req.GetBaseModelID()),
	})
	if err != nil {
		return nil, fmt.Errorf("s.commonDB.GetTrainedModels: %w", err)
	}

	resp := &pb.GetTrainedModelsResponse{
		Models: make([]*pb.ShortTrainedModel, 0, len(modelRows)),
		PageInfo: &pb.PageInfo{
			Offset: req.GetOffset(),
			Limit:  req.GetLimit(),
		},
	}

	if len(modelRows) == 0 {
		return resp, nil
	}
	resp.PageInfo.Total = uint64(modelRows[0].Count.Int64)

	for _, row := range modelRows {
		resp.Models = append(resp.Models, &pb.ShortTrainedModel{
			TrainedModelID:   uint64(row.ID),
			Name:             row.Name,
			TrainStatus:      convertTrainingStatus(row.LaunchStatus),
			BaseModelID:      uint64(row.ModelID),
			BaseModelName:    row.ModelName,
			ProblemName:      row.ProblemName,
			TrainDatasetID:   uint64(row.TrainingDatasetID),
			TrainDatasetName: row.TrainingDatasetName,
			CreatedAt:        timestamppb.New(row.CreatedAt.Time),
			LaunchID:         uint64(row.LaunchID.Int64),
		})
	}

	return resp, nil
}

func convertTrainingStatus(status string) pb.LaunchStatus {
	switch status {
	case pb.LaunchStatus_LaunchStatusNotStarted.String():
		return pb.LaunchStatus_LaunchStatusNotStarted
	case pb.LaunchStatus_LaunchStatusInProgress.String():
		return pb.LaunchStatus_LaunchStatusInProgress
	case pb.LaunchStatus_LaunchStatusError.String():
		return pb.LaunchStatus_LaunchStatusError
	case pb.LaunchStatus_LaunchStatusSuccess.String():
		return pb.LaunchStatus_LaunchStatusSuccess
	default:
		return pb.LaunchStatus_LaunchStatusUnknown
	}
}
