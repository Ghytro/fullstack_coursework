package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/sirupsen/logrus"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
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

// RunInTransaction запустить коллбек в транзакции либо взятой из контекста либо в новой.
// Если контекст в транзакции был, верим "наслово", что транзакция была проинициализирована
// с контекстом уже на слое репозиториев. Также гарантируется и успешный коммит/роллбек транзакции
// на верхних слоях
func (db *PGDB) RunInTransaction(ctx context.Context, f func(*TX) error) error {
	return db.WithContext(ctx).RunInTransaction(ctx, func(tx *TX) error {
		return f(tx)
	})
}

// WithContext замена базы данных транзакцией из контекста.
// Если транзакции не было применяем контекст к бд
func (db *PGDB) WithContext(ctx context.Context) DBI {
	if tx := GetTx(ctx); tx != nil {
		return tx
	}
	return &PGDB{
		DB: db.DB.WithContext(ctx),
	}
}

// BeginContext начать транзакцию и вернуть ее объект
func (db *PGDB) BeginContext(ctx context.Context) (*TX, error) {
	tx, err := db.DB.BeginContext(ctx)
	if err != nil {
		return nil, err
	}
	return &TX{
		Tx: tx,
	}, nil
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
	if tx := GetTx(ctx); tx != nil { // странный кейс но все же
		return f(tx)
	}
	// todo: может спавнить сейвпоинты для вложенных транзакций?
	return f(t)
}

func (t *TX) BeginContext(ctx context.Context) (*TX, error) {
	return t, nil
}

func (t *TX) WithContext(ctx context.Context) DBI {
	logrus.StandardLogger().Warn("parent tx context initialization not allowed")
	return t
}

type DBI interface {
	Exec(query interface{}, params ...interface{}) (orm.Result, error)
	ExecContext(ctx context.Context, query interface{}, params ...interface{}) (orm.Result, error)

	Query(model, query interface{}, params ...interface{}) (orm.Result, error)
	QueryContext(ctx context.Context, model, query interface{}, params ...interface{}) (orm.Result, error)

	Model(model ...interface{}) *orm.Query
	ModelContext(ctx context.Context, model ...interface{}) *orm.Query

	RunInTransaction(ctx context.Context, fn func(*TX) error) error

	BeginContext(ctx context.Context) (*TX, error)

	WithContext(ctx context.Context) DBI
}
