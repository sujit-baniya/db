package db

import (
	"context"
	"sync"
)

type Entity any

type Field struct {
	Column string
	Value  any
}

type repositories struct {
	repos map[string]any
	mu    *sync.RWMutex
}

func (r *repositories) Add(key string, repo any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.repos[key] = repo
}

func (r *repositories) Get(key string) (any, bool) {
	val, ok := r.repos[key]
	return val, ok
}

var rs *repositories

func init() {
	rs = &repositories{mu: &sync.RWMutex{}, repos: make(map[string]any)}
}

type Repository[T Entity] interface {
	Find(ctx context.Context, id any) (*T, error)
	FindBy(ctx context.Context, fs ...Field) ([]T, error)
	All(ctx context.Context) ([]T, error)
	FindByWithRelations(ctx context.Context, fs ...Field) ([]T, error)
	FindWithRelations(ctx context.Context, id any) (*T, error)
	FindFirstBy(ctx context.Context, fs ...Field) (*T, error)
	Create(ctx context.Context, t *T) error
	Raw(ctx context.Context, sql string, values ...any) ([]T, error)
	RawAny(ctx context.Context, rs any, sql string, values ...any) (any, error)
	RawMapFirst(ctx context.Context, sql string, values ...any) (map[string]any, error)
	RawMapSlice(ctx context.Context, sql string, values ...any) ([]map[string]any, error)
	CreateBulk(ctx context.Context, ts []T) error
	Update(ctx context.Context, t *T, fs ...Field) error
	Delete(ctx context.Context, t *T) error
	GetDB() any
}
