package models

import (
	"context"

	pb "github.com/hse-experiments-platform/models/pkg/models"
)

func (s *modelsService) GetTrainedModels(ctx context.Context, req *pb.GetTrainedModelsRequest) (*pb.GetTrainedModelsResponse, error) {
	return &pb.GetTrainedModelsResponse{
		Models: []*pb.ShortTrainedModel{
			{
				TrainedModelID:   1,
				Name:             "Trained linear regression",
				TrainStatus:      pb.TrainStatus_TrainStatusDone,
				BaseModelID:      1,
				BaseModelName:    "Линейная регрессия",
				ProblemName:      "Регрессия",
				TrainDatasetID:   2,
				TrainDatasetName: "Mock dataset",
				LaunchID:         123,
			},
			{
				TrainedModelID:   2,
				Name:             "Trained logistic regression",
				TrainStatus:      pb.TrainStatus_TrainStatusError,
				BaseModelID:      2,
				BaseModelName:    "Логистическая регрессия",
				ProblemName:      "Классификация",
				TrainDatasetID:   4,
				TrainDatasetName: "Mock dataset 2",
				LaunchID:         666,
			},
		},
		PageInfo: &pb.PageInfo{
			Offset: req.GetOffset(),
			Limit:  req.GetLimit(),
			Total:  2,
		},
	}, nil
}
