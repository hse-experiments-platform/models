// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"
)

type Querier interface {
	GetModel(ctx context.Context, id int64) (GetModelRow, error)
	GetModelHyperparameters(ctx context.Context, id int64) ([]GetModelHyperparametersRow, error)
	GetModelProblem(ctx context.Context, id int64) (GetModelProblemRow, error)
	GetModels(ctx context.Context, arg GetModelsParams) ([]GetModelsRow, error)
	GetProblems(ctx context.Context, arg GetProblemsParams) ([]GetProblemsRow, error)
}

var _ Querier = (*Queries)(nil)
