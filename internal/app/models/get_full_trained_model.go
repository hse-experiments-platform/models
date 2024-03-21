package models

import (
	"context"

	pb "github.com/hse-experiments-platform/models/pkg/models"
)

func (s *modelsService) GetFullTrainedModel(ctx context.Context, req *pb.GetFullTrainedModelRequest) (*pb.GetFullTrainedModelResponse, error) {
	if req.GetTrainedModelID() == 1 {
		return &pb.GetFullTrainedModelResponse{Model: &pb.TrainedModel{
			TrainedModelID: 1,
			Name:           "Trained linear regression",
			TrainStatus:    pb.TrainStatus_TrainStatusError,
			BaseModelID:    1,
			BaseModelName:  "Линейная регрессия",
			Problem: &pb.ShortProblem{
				Id:          1,
				Name:        "Регрессия",
				Description: "Задача, в которой модель будет побирать вещественное число в качестве ответа",
			},
			TrainDatasetID:   2,
			TrainDatasetName: "Mock dataset",
			TargetColumn:     "target_mock",
			Schema: &pb.DatasetSchema{Columns: []*pb.DatasetSchema_SchemaColumn{{
				Name: "target_mock",
				Type: "int",
			}}},
			LaunchID: 123,
		}}, nil
	} else {
		return &pb.GetFullTrainedModelResponse{Model: &pb.TrainedModel{
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
			TargetColumn:     "y",
			Schema: &pb.DatasetSchema{Columns: []*pb.DatasetSchema_SchemaColumn{
				{
					Name: "x",
					Type: "int",
				},
				{
					Name: "y",
					Type: "int",
				},
			}},
		}}, nil
	}
}
