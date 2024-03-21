// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type DatasetStatus string

const (
	DatasetStatusInitializing DatasetStatus = "initializing"
	DatasetStatusLoading      DatasetStatus = "loading"
	DatasetStatusReady        DatasetStatus = "ready"
	DatasetStatusError        DatasetStatus = "error"
)

func (e *DatasetStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = DatasetStatus(s)
	case string:
		*e = DatasetStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for DatasetStatus: %T", src)
	}
	return nil
}

type NullDatasetStatus struct {
	DatasetStatus DatasetStatus
	Valid         bool // Valid is true if DatasetStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullDatasetStatus) Scan(value interface{}) error {
	if value == nil {
		ns.DatasetStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.DatasetStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullDatasetStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.DatasetStatus), nil
}

type LaunchType string

const (
	LaunchTypeTrain   LaunchType = "train"
	LaunchTypePredict LaunchType = "predict"
	LaunchTypeGeneric LaunchType = "generic"
)

func (e *LaunchType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = LaunchType(s)
	case string:
		*e = LaunchType(s)
	default:
		return fmt.Errorf("unsupported scan type for LaunchType: %T", src)
	}
	return nil
}

type NullLaunchType struct {
	LaunchType LaunchType
	Valid      bool // Valid is true if LaunchType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullLaunchType) Scan(value interface{}) error {
	if value == nil {
		ns.LaunchType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.LaunchType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullLaunchType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.LaunchType), nil
}

type ModelTrainingStatus string

const (
	ModelTrainingStatusNotStarted ModelTrainingStatus = "not_started"
	ModelTrainingStatusInProgress ModelTrainingStatus = "in_progress"
	ModelTrainingStatusError      ModelTrainingStatus = "error"
	ModelTrainingStatusDone       ModelTrainingStatus = "done"
)

func (e *ModelTrainingStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ModelTrainingStatus(s)
	case string:
		*e = ModelTrainingStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ModelTrainingStatus: %T", src)
	}
	return nil
}

type NullModelTrainingStatus struct {
	ModelTrainingStatus ModelTrainingStatus
	Valid               bool // Valid is true if ModelTrainingStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullModelTrainingStatus) Scan(value interface{}) error {
	if value == nil {
		ns.ModelTrainingStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ModelTrainingStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullModelTrainingStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ModelTrainingStatus), nil
}

type Dataset struct {
	ID          int64
	Name        string
	Version     string
	Status      DatasetStatus
	RowsCount   int64
	CreatorID   int64
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	UploadError pgtype.Text
}

type DatasetSchema struct {
	DatasetID    int64
	ColumnNumber int32
	ColumnName   string
	ColumnType   string
}

type Hyperparameter struct {
	ID           int64
	Name         string
	Description  string
	Type         string
	DefaultValue []byte
	ModelID      pgtype.Int8
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

type Launch struct {
	ID          int64
	LaunchType  LaunchType
	Name        string
	Description string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	FinishedAt  pgtype.Timestamptz
	LaunchError pgtype.Text
}

type Metric struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type Model struct {
	ID          int64
	Name        string
	Description string
	ProblemID   int64
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type PredictResult struct {
	LaunchID        int64
	TrainedModelID  int64
	InputDatasetID  int64
	Status          DatasetStatus
	OutputDatasetID int64
}

type Problem struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type ProblemMetric struct {
	ProblemID int64
	MetricID  int64
}

type TrainMetric struct {
	LaunchID       int64
	TrainedModelID int64
	MetricID       int64
	Value          []byte
}

type TrainedModel struct {
	ID                  int64
	Name                string
	Description         string
	ModelID             int64
	ModelTrainingStatus ModelTrainingStatus
	TrainingDatasetID   int64
	TargetColumn        string
	TrainError          pgtype.Text
	CreatedAt           pgtype.Timestamptz
	UpdatedAt           pgtype.Timestamptz
	LaunchID            int64
}

type User struct {
	ID        int64
	GoogleID  string
	Email     string
	CreatedAt pgtype.Timestamptz
}
