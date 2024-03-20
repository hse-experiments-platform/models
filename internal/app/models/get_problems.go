package models

import (
	"context"
	"fmt"

	"github.com/hse-experiments-platform/models/internal/pkg/storage/db"
	pb "github.com/hse-experiments-platform/models/pkg/models"
)

func (s *modelsService) GetProblems(ctx context.Context, req *pb.GetProblemsRequest) (*pb.GetProblemsResponse, error) {
	problem, err := s.commonDB.GetProblems(ctx, db.GetProblemsParams{
		Name:   "%" + req.GetQuery() + "%",
		Limit:  int64(req.GetLimit()),
		Offset: int64(req.GetOffset()),
	})
	if err != nil {
		return nil, fmt.Errorf("s.commonDB.GetModels: %w", err)
	}

	resp := &pb.GetProblemsResponse{
		Problems: make([]*pb.ShortProblem, 0, len(problem)),
		PageInfo: &pb.PageInfo{
			Offset: req.GetOffset(),
			Limit:  req.GetLimit(),
		},
	}

	if len(problem) == 0 {
		return resp, nil
	}
	resp.PageInfo.Total = uint64(problem[0].Count)

	for _, row := range problem {
		resp.Problems = append(resp.Problems, &pb.ShortProblem{
			Id:          uint64(row.ID),
			Name:        row.Name,
			Description: row.Description,
		})
	}

	return resp, nil
}
