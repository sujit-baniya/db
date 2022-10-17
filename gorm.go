package db

import (
	"context"
	"errors"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Gorm[T Entity] struct {
	db *gorm.DB
}

func NewGormRepository[T Entity](ctx context.Context, db *gorm.DB) Repository[T] {
	name := reflect.TypeOf((*T)(nil)).Elem().Name()
	if r, ok := rs.Get(name); ok {
		return r.(Repository[T])
	}
	repo := Gorm[T]{db}
	rs.Add(name, repo)
	return repo
}

func (r Gorm[T]) Find(ctx context.Context, id any) (*T, error) {
	var t T

	result := r.db.WithContext(ctx).Where(id).First(&t)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &t, nil
}
func (r Gorm[T]) FindBy(ctx context.Context, fs ...Field) ([]T, error) {
	var t []T

	whereClause := make(map[string]interface{}, len(fs))

	for _, f := range fs {
		whereClause[f.Column] = f.Value
	}

	result := r.db.WithContext(ctx).Where(whereClause).Find(&t)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return t, nil
}
func (r Gorm[T]) All(ctx context.Context) ([]T, error) {
	var t []T
	result := r.db.WithContext(ctx).Find(&t)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return t, nil
}
func (r Gorm[T]) FindByWithRelations(ctx context.Context, fs ...Field) ([]T, error) {
	var t []T

	whereClause := make(map[string]interface{}, len(fs))

	for _, f := range fs {
		whereClause[f.Column] = f.Value
	}

	result := r.db.WithContext(ctx).Preload(clause.Associations).Where(whereClause).Find(&t)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return t, nil
}
func (r Gorm[T]) FindWithRelations(ctx context.Context, id any) (*T, error) {
	var t T

	result := r.db.WithContext(ctx).Preload(clause.Associations).Where(id).First(&t)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &t, nil
}
func (r Gorm[T]) FindFirstBy(ctx context.Context, fs ...Field) (*T, error) {
	ts, err := r.FindBy(ctx, fs...)
	if err != nil {
		return nil, err
	}

	if len(ts) >= 1 {
		return &ts[0], nil
	}

	return nil, errors.New("Record not found")
}
func (r Gorm[T]) Create(ctx context.Context, t *T) error {
	return r.db.WithContext(ctx).Create(t).Error
}
func (r Gorm[T]) Raw(ctx context.Context, sql string, values ...any) ([]T, error) {
	var ts []T
	result := r.db.WithContext(ctx).Raw(sql, values...).Scan(&ts)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return ts, nil
}
func (r Gorm[T]) RawAny(ctx context.Context, rs any, sql string, values ...any) (any, error) {
	result := r.db.WithContext(ctx).Raw(sql, values...).Scan(&rs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return rs, nil
}
func (r Gorm[T]) RawMapFirst(ctx context.Context, sql string, values ...any) (map[string]any, error) {
	var rt map[string]any
	result := r.db.WithContext(ctx).Raw(sql, values...).Scan(&rt)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return rt, nil
}
func (r Gorm[T]) RawMapSlice(ctx context.Context, sql string, values ...any) ([]map[string]any, error) {
	var rt []map[string]any
	result := r.db.WithContext(ctx).Raw(sql, values...).Scan(&rt)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return rt, nil
}
func (r Gorm[T]) CreateBulk(ctx context.Context, ts []T) error {
	return r.db.WithContext(ctx).Create(&ts).Error
}
func (r Gorm[T]) Update(ctx context.Context, t *T, fs ...Field) error {
	updateFields := make(map[string]interface{}, len(fs))

	for _, f := range fs {
		updateFields[f.Column] = f.Value
	}

	return r.db.WithContext(ctx).Model(t).Updates(updateFields).Error
}
func (r Gorm[T]) Delete(ctx context.Context, t *T) error {
	return r.db.WithContext(ctx).Delete(t).Error
}
func (r Gorm[T]) GetDB() any {
	return r.db
}
