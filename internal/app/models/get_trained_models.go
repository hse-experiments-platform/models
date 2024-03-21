package models

import (
	"context"

	pb "github.com/hse-experiments-platform/models/pkg/models"
)

func (s *modelsService) GetTrainedModels(ctx context.Context, req *pb.GetTrainedModelsRequest) (*pb.GetTrainedModelsResponse, error) {
	return &pb.GetTrainedModelsResponse{
		Models: []*pb.ShortTrainedModel{
			{
				TrainedModelID: 1,
				Name:           "Trained linear regression",
				TrainStatus:    pb.TrainStatus_TrainStatusDone,
				BaseModelID:    1,
				BaseModelName:  "Линейная регрессия",
				Problem: &pb.ShortProblem{
					Id:          1,
					Name:        "Регрессия",
					Description: "Задача, в которой модель будет побирать вещественное число в качестве ответа",
				},
				TrainDatasetID:   2,
				TrainDatasetName: "Mock dataset",
				LaunchID:         123,
			},
			{
				TrainedModelID: 2,
				Name:           "Trained logistic regression",
				TrainStatus:    pb.TrainStatus_TrainStatusError,
				BaseModelID:    2,
				BaseModelName:  "Логистическая регрессия",
				Problem: &pb.ShortProblem{
					Id:          2,
					Name:        "Классификация",
					Description: "Задача, в которой модель будет распределять объекты по заданным классам",
				},
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
