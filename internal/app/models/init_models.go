package models

import (
	"context"

	"github.com/hse-experiments-platform/models/internal/pkg/core"
	"github.com/hse-experiments-platform/models/internal/pkg/storage/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

func convertRow(row db.GetAllModelsRow) *core.ModelConfig {
	hyperparameters := make([]core.Hyperparameter, 0, len(row.HyperparameterNames))

	for i, hyperparameter := range row.HyperparameterNames {
		hyperparameters = append(hyperparameters, core.Hyperparameter{
			Name:         hyperparameter,
			Description:  row.HyperparameterDescriptions[i],
			Type:         row.HyperparameterTypes[i],
			DefaultValue: string(row.HyperparameterDefaultValues[i]),
		})
	}

	return &core.ModelConfig{
		ModelName:            row.Name,
		ModelDescription:     row.Description,
		ClassName:            row.ClassName.String,
		Problem:              row.ProblemName,
		TrainMetrics:         row.MetricNames,
		TrainHyperparameters: hyperparameters,
	}
}
func (s *modelsService) insertModels(ctx context.Context, modelHashes map[string]*core.ModelConfig) func(tx pgx.Tx) error {
	return func(tx pgx.Tx) (err error) {
		defer func() {
			if err != nil {
				err = tx.Rollback(ctx)
				return
			}
		}()

		txDB := s.commonDB.WithTx(tx)

		rows, err := txDB.GetAllModels(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get models from the database")
			return err
		}
		dbHashes := make(map[string]*core.ModelConfig)
		for _, row := range rows {
			m := convertRow(row)
			hash, err := m.Hash()
			if err != nil {
				log.Error().Err(err).Msg("Failed to compute model hash")
				return err
			}
			dbHashes[hash] = m
		}

		log.Debug().Any("modelHashes", modelHashes).Any("dbHashes", dbHashes).Msg("test")

		for h, m := range modelHashes {
			if _, ok := dbHashes[h]; !ok {
				log.Info().Msgf("Inserting model %s", m.ModelName)
				id, err := txDB.CreateModel(ctx, db.CreateModelParams{
					Name:        pgtype.Text{String: m.ModelName, Valid: true},
					Description: pgtype.Text{String: m.ModelDescription, Valid: true},
					Problem:     pgtype.Text{String: m.Problem, Valid: true},
					ClassName:   pgtype.Text{String: m.ClassName, Valid: true},
				})
				if err != nil {
					log.Error().Err(err).Msg("Failed to insert model")
					return err
				}

				err = txDB.CreateModelMetrics(ctx, db.CreateModelMetricsParams{
					MetricNames: m.TrainMetrics,
					ModelID:     id,
				})
				if err != nil {
					log.Error().Err(err).Msg("Failed to insert model metrics")
					return err
				}

				names := make([]string, 0, len(m.TrainHyperparameters))
				descriptions := make([]string, 0, len(m.TrainHyperparameters))
				types := make([]string, 0, len(m.TrainHyperparameters))
				defaultValues := make([]string, 0, len(m.TrainHyperparameters))
				for _, hyperparameter := range m.TrainHyperparameters {
					names = append(names, hyperparameter.Name)
					descriptions = append(descriptions, hyperparameter.Description)
					types = append(types, hyperparameter.Type)
					defaultValues = append(defaultValues, hyperparameter.DefaultValue)
				}

				err = txDB.CreateHyperparameters(ctx, db.CreateHyperparametersParams{
					Names:         names,
					Descriptions:  descriptions,
					Types:         types,
					DefaultValues: defaultValues,
					ModelID:       pgtype.Int8{Int64: id, Valid: true},
				})
				if err != nil {
					log.Error().Err(err).Msg("Failed to insert hyperparameters")
					return err
				}
			} else {
				log.Info().Msgf("Model %s is already in db", m.ModelName)
			}
		}

		return nil
	}
}
func (s *modelsService) InitModels(ctx context.Context, dirPath string) error {
	log.Info().Msgf("Initializing models from %s", dirPath)

	models, err := core.ParseConfigFromSubfolders(dirPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse models")
		return err
	}
	modelHashes := make(map[string]*core.ModelConfig)
	for _, model := range models {
		hash, err := model.Hash()
		if err != nil {
			log.Error().Err(err).Msg("Failed to compute model hash")
			return err
		}
		modelHashes[hash] = model
	}

	err = pgx.BeginTxFunc(ctx, s.commonDBConn, pgx.TxOptions{AccessMode: pgx.ReadWrite}, s.insertModels(ctx, modelHashes))
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert models")
		return err
	}

	return nil
}
