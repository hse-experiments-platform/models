package models

import (
	"context"

	"github.com/hse-experiments-platform/library/pkg/utils/token"
	"github.com/hse-experiments-platform/models/internal/pkg/storage/db"
	"github.com/hse-experiments-platform/models/internal/pkg/storage/mlflowdb"
	pb "github.com/hse-experiments-platform/models/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.ModelsServiceServer = (*modelsService)(nil)

type modelsService struct {
	pb.UnimplementedModelsServiceServer
	commonDBConn *pgxpool.Pool
	mlflowDBConn *pgxpool.Pool
	maker        token.Maker

	commonDB *db.Queries
	mlflowDB *mlflowdb.Queries
}

func NewService(commonDBConn *pgxpool.Pool, mlflowDB *pgxpool.Pool, maker token.Maker) *modelsService {
	return &modelsService{
		commonDBConn: commonDBConn,
		maker:        maker,

		commonDB: db.New(commonDBConn),
		mlflowDB: mlflowdb.New(mlflowDB),
	}
}

func getUserID(ctx context.Context) (int64, error) {
	var userID int64
	userID, ok := ctx.Value(token.UserIDContextKey).(int64)
	if !ok {
		log.Error().Msgf("invalid userID context key type: %T", ctx.Value(token.UserIDContextKey))
		return 0, status.New(codes.Internal, "internal error").Err()
	}

	return userID, nil
}
