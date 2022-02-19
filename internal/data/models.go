package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Payments PaymentsModel
	Categories CategoriesModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Payments: PaymentsModel{DB: db},
		Categories: CategoriesModel{DB: db},
	}
}
