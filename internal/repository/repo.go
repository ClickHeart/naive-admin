package repository

import (
	"context"
	"naive-admin/pkg/db"

	"gorm.io/gorm"
)

type contextKey string

const TxKey contextKey = "TxKey"

type Repository struct{}

var Repo = &Repository{}

// DB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (r *Repository) DB(c context.Context) *gorm.DB {
	v := c.Value(TxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return db.Dao.WithContext(c)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return db.Dao.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, TxKey, tx)
		return fn(ctx)
	})
}
