package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Ghytro/galleryapp/internal/common"

	"github.com/go-pg/pg/v10"
)

type PGDB struct {
	*pg.DB
}

func NewPGDB(url string, queryLogger *PGLogger) *PGDB {
	dbOpts, err := pg.ParseURL(os.Getenv("DB_URL"))
	common.LogFatalErr(err)
	db := pg.Connect(dbOpts)
	db.AddQueryHook(queryLogger)
	return &PGDB{DB: db}
}

func (db *PGDB) RunInTransaction(ctx context.Context, f func(*TX) error) error {
	return db.DB.RunInTransaction(ctx, func(tx *pg.Tx) error {
		return f(&TX{Tx: tx})
	})
}

type PGLogger struct {
}

func (l *PGLogger) BeforeQuery(ctx context.Context, e *pg.QueryEvent) (context.Context, error) {
	q, err := e.FormattedQuery()
	if err != nil {
		log.Println(err) // todo: нормальный логгер
	}
	log.Println(string(q))
	return ctx, nil
}

func (l *PGLogger) AfterQuery(ctx context.Context, e *pg.QueryEvent) error {
	log.Printf("executed at: %dms", time.Since(e.StartTime).Milliseconds())
	return nil
}

type TX struct {
	*pg.Tx
}

func (t *TX) RunInTransaction(ctx context.Context, f func(*TX) error) error {
	// todo: может спавнить сейвпоинты для вложенных транзакций?
	return f(t)
}
