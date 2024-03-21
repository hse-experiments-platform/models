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
			TrainStatus:      convertTrainingStatus(row.ModelTrainingStatus),
			BaseModelID:      uint64(row.ModelID),
			BaseModelName:    row.ModelName,
			ProblemName:      row.ProblemName,
			TrainDatasetID:   uint64(row.TrainingDatasetID),
			TrainDatasetName: row.TrainingDatasetName,
			CreatedAt:        timestamppb.New(row.CreatedAt.Time),
			LaunchID:         uint64(row.LaunchID),
		})
	}

	return resp, nil
}

func convertTrainingStatus(s db.ModelTrainingStatus) pb.TrainStatus {
	switch s {
	case db.ModelTrainingStatusNotStarted:
		return pb.TrainStatus_TrainStatusNotStarted
	case db.ModelTrainingStatusInProgress:
		return pb.TrainStatus_TrainStatusInProgress
	case db.ModelTrainingStatusError:
		return pb.TrainStatus_TrainStatusError
	case db.ModelTrainingStatusDone:
		return pb.TrainStatus_TrainStatusDone
	default:
		return pb.TrainStatus_TrainStatusUnknown
	}
}
