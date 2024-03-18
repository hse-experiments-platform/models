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

type Dataset struct {
	ID        int64
	Name      string
	Version   string
	Status    DatasetStatus
	RowsCount int64
	CreatorID int64
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type User struct {
	ID        int64
	GoogleID  string
	Email     string
	CreatedAt pgtype.Timestamptz
}
