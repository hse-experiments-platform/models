package models

import (
	"context"
	"fmt"

	"github.com/hse-experiments-platform/models/internal/pkg/storage/db"
	pb "github.com/hse-experiments-platform/models/pkg/models"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *modelsService) GetPredictions(ctx context.Context, req *pb.GetPredictionsRequest) (*pb.GetPredictionsResponse, error) {
	predictionRows, err := s.commonDB.GetPredictions(ctx, db.GetPredictionsParams{
		LaunchStatus:   pb.LaunchStatus_LaunchStatusSuccess.String(),
		LaunchType:     "LaunchTypePredict",
		Limit:          int64(req.GetLimit()),
		Offset:         int64(req.GetOffset()),
		TrainedModelID: int64(req.GetTrainedModelID()),
		Name:           "%" + req.GetQuery() + "%",
	})
	log.Debug().Any("df", db.GetPredictionsParams{
		LaunchStatus:   pb.LaunchStatus_LaunchStatusSuccess.String(),
		LaunchType:     "LaunchTypePredict",
		Limit:          int64(req.GetLimit()),
		Offset:         int64(req.GetOffset()),
		TrainedModelID: int64(req.GetTrainedModelID()),
		Name:           "%" + req.GetQuery() + "%",
	}).Msg("q23")
	if err != nil {
		return nil, fmt.Errorf("s.commonDB.GetPredictions: %w", err)
	}

	resp := &pb.GetPredictionsResponse{
		Predictions: make([]*pb.PredictionInfo, 0, len(predictionRows)),
		PageInfo: &pb.PageInfo{
			Offset: req.GetOffset(),
			Limit:  req.GetLimit(),
		},
	}

	if len(predictionRows) == 0 {
		return resp, nil
	}
	resp.PageInfo.Total = uint64(predictionRows[0].Count.Int64)

	for _, row := range predictionRows {
		resp.Predictions = append(resp.Predictions, &pb.PredictionInfo{
			LaunchID:      row.ID,
			Name:          row.Name,
			Status:        convertTrainingStatus(row.LaunchStatus),
			DatasetName:   row.DatasetName,
			Target:        row.TargetCol,
			StartDateTime: timestamppb.New(row.CreatedAt.Time),
		})
	}

	return resp, nil
}
