package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/hse-experiments-platform/models/internal/pkg/storage/db"
	pb "github.com/hse-experiments-platform/models/pkg/models"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *modelsService) GetFullModel(ctx context.Context, req *pb.GetFullModelRequest) (*pb.GetFullModelResponse, error) {
	resp := &pb.GetFullModelResponse{Model: &pb.Model{}}

	err := pgx.BeginTxFunc(ctx, s.commonDBConn, pgx.TxOptions{AccessMode: pgx.ReadOnly}, func(tx pgx.Tx) error {
		txdb := s.commonDB.WithTx(tx)
		modelID := int64(req.GetModelID())

		model, err := txdb.GetModel(ctx, modelID)
		if errors.Is(err, pgx.ErrNoRows) {
			return status.Error(codes.NotFound, fmt.Sprintf("model with id %v not found", modelID))
		} else if err != nil {
			return fmt.Errorf("txdb.GetModel: %w", err)
		}

		params, err := txdb.GetModelHyperparameters(ctx, modelID)
		if err != nil {
			return fmt.Errorf("txdb.GetModelHyperparameters: %w", err)
		}

		pr, err := txdb.GetModelProblem(ctx, modelID)
		if errors.Is(err, pgx.ErrNoRows) {
			return status.Error(codes.NotFound, fmt.Sprintf("problem for model with id %v not found", modelID))
		} else if len(pr.MetricIds) != len(pr.MetricDescriptions) || len(pr.MetricDescriptions) != len(pr.MetricNames) {
			return status.Errorf(codes.Internal, "invalid arrays len in problem with id %v", pr.ID)
		} else if err != nil {
			return fmt.Errorf("txdb.GetModelProblem: %w", err)
		}

		resp.Model = &pb.Model{
			ModelId:         uint64(modelID),
			Name:            model.Name,
			Description:     model.Description,
			Hyperparameters: convertHyperparameters(params),
			Problem:         convertProblem(pr),
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("pgx.BeginTxFunc: %w", err)
	}

	return resp, nil
}

func convertHyperparameters(params []db.GetModelHyperparametersRow) []*pb.Hyperparameter {
	res := make([]*pb.Hyperparameter, 0, len(params))

	for _, p := range params {
		res = append(res, &pb.Hyperparameter{
			Id:           uint64(p.ID),
			Name:         p.Name,
			Description:  p.Description,
			Type:         p.Type,
			DefaultValue: structpb.NewStringValue(string(p.DefaultValue)),
		})
	}

	return res
}

func convertProblem(p db.GetModelProblemRow) *pb.Problem {
	resp := &pb.Problem{
		Id:          uint64(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Metrics:     make([]*pb.Metric, 0, len(p.MetricIds)),
	}

	for i, _ := range p.MetricIds {
		resp.Metrics = append(resp.Metrics, &pb.Metric{
			Id:          uint64(p.MetricIds[i]),
			Name:        p.MetricNames[i],
			Description: p.MetricDescriptions[i],
		})
	}

	return resp
}
