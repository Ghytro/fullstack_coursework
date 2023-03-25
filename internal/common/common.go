package common

import (
	"fmt"
	"log"
	"sync"

	"github.com/biter777/countries"
	"github.com/go-pg/pg/v10/orm"
)

func LogFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type SyncMap[K comparable, V any] struct {
	mut *sync.Mutex
	m   map[K]V
}

func NewSyncMap[K comparable, V any](mutex *sync.Mutex) *SyncMap[K, V] {
	return &SyncMap[K, V]{
		mut: mutex,
		m:   make(map[K]V),
	}
}

func (m *SyncMap[K, V]) Get(key K) (V, bool) {
	m.mut.Lock()
	defer m.mut.Unlock()
	val, ok := m.m[key]
	return val, ok
}

func (m *SyncMap[K, V]) MustGet(key K) V {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.m[key]
}

func (m *SyncMap[K, V]) Set(key K, val V) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.m[key] = val
}

func GetCountryByAlpha2(alpha2 string) *countries.Country {
	for _, c := range countries.AllInfo() {
		if c.Alpha2 == alpha2 {
			return c
		}
	}
	return nil
}

type Range[T any] struct {
	From, To T
}

type PageData struct {
	Page, PageSize int
}

type StringDataFilter struct {
	Value    string
	Likeness StrDataLikeness
}

type StrDataLikeness uint8

const (
	StrDataLikenessExact = iota
	StrDataLikenessSubstr
)

func (l StringDataFilter) Apply(q *orm.Query, columnName string) {
	switch l.Likeness {
	case StrDataLikenessExact:
		q.Where(fmt.Sprintf("%s = ?", columnName), l.Value)
	case StrDataLikenessSubstr:
		q.Where(fmt.Sprintf("position(%s in ?) > 0", columnName), l.Value)
	}
}
