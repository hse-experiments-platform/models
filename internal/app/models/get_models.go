package models

import (
	"context"
	"fmt"

	"github.com/hse-experiments-platform/models/internal/pkg/storage/db"
	pb "github.com/hse-experiments-platform/models/pkg/models"
)

func (s *modelsService) GetModels(ctx context.Context, req *pb.GetModelsRequest) (*pb.GetModelsResponse, error) {
	modelRows, err := s.commonDB.GetModels(ctx, db.GetModelsParams{
		Name:      "%" + req.GetQuery() + "%",
		Limit:     int64(req.GetLimit()),
		Offset:    int64(req.GetOffset()),
		ProblemID: int64(req.GetProblemID()),
	})
	if err != nil {
		return nil, fmt.Errorf("s.commonDB.GetModels: %w", err)
	}

	resp := &pb.GetModelsResponse{
		Models: make([]*pb.ShortModel, 0, len(modelRows)),
		PageInfo: &pb.PageInfo{
			Offset: req.GetOffset(),
			Limit:  req.GetLimit(),
		},
	}

	if len(modelRows) == 0 {
		return resp, nil
	}
	resp.PageInfo.Total = uint64(modelRows[0].Count)

	for _, row := range modelRows {
		resp.Models = append(resp.Models, &pb.ShortModel{
			ModelId:     uint64(row.ID),
			Name:        row.Name,
			Description: row.Description,
		})
	}

	return resp, nil
}
