package repository

import (
	"context"

	"github.com/Ghytro/galleryapp/internal/database"

	"github.com/go-pg/pg/v10/orm"
)

type DBI interface {
	Exec(query interface{}, params ...interface{}) (orm.Result, error)
	ExecContext(ctx context.Context, query interface{}, params ...interface{}) (orm.Result, error)

	Query(repository, query interface{}, params ...interface{}) (orm.Result, error)
	QueryContext(ctx context.Context, repository, query interface{}, params ...interface{}) (orm.Result, error)

	Model(repository ...interface{}) *orm.Query
	ModelContext(ctx context.Context, repository ...interface{}) *orm.Query

	RunInTransaction(ctx context.Context, fn func(*database.TX) error) error
}
