package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	pb "github.com/hse-experiments-platform/models/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *modelsService) GetTrainMetrics(ctx context.Context, req *pb.GetTrainMetricsRequest) (*pb.GetTrainMetricsResponse, error) {
	runIDJson, err := s.commonDB.GetTrainedModelRunID(ctx, int64(req.GetTrainedModelID()))
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "trained model not found")
	} else if err != nil {
		return nil, fmt.Errorf("s.commonDB.GetTrainedModelRunID: %w", err)
	}

	var runID string
	if err := json.Unmarshal(runIDJson, &runID); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	metrics, err := s.mlflowDB.GetMetrics(ctx, pgtype.Text{String: string(runID), Valid: true})
	if err != nil {
		return nil, fmt.Errorf("s.mlflowDB.GetMetrics: %w", err)
	}

	if len(metrics) == 0 {
		return &pb.GetTrainMetricsResponse{}, nil
	}

	maxStep := metrics[len(metrics)-1].Step

	resp := &pb.GetTrainMetricsResponse{
		Metrics:   make(map[string]float64),
		CvMetrics: make([]*pb.CVMetrics, maxStep),
	}

	for i := 0; i < int(maxStep); i++ {
		resp.CvMetrics[i] = &pb.CVMetrics{
			Metrics: make(map[string]float64),
		}
	}

	for _, metric := range metrics {
		if metric.Step == maxStep {
			resp.Metrics[metric.Key] = metric.Value
		} else {
			resp.CvMetrics[metric.Step].Metrics[metric.Key] = metric.Value
		}
	}

	return resp, nil
}
