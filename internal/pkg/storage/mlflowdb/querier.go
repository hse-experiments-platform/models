// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package mlflowdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	GetMetrics(ctx context.Context, runID pgtype.Text) ([]GetMetricsRow, error)
}

var _ Querier = (*Queries)(nil)
