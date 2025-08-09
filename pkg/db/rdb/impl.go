package rdb

import (
	"context"
	"reflect"
	"strings"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Engine 泛型仓库实现
type Engine[T any] struct {
	db *gorm.DB
}

// NewEngine 创建新的泛型仓库
func NewEngine[T any](db *gorm.DB) EngineInterface[T] {
	return &Engine[T]{db: db}
}

// DB 返回原始的 gorm.DB 实例
func (e *Engine[T]) DB() *gorm.DB {
	return e.db
}

func (e *Engine[T]) Create(ctx context.Context, entity *T) error {
	return e.db.WithContext(ctx).Create(entity).Error
}

func (e *Engine[T]) GetByID(ctx context.Context, id any) (*T, error) {
	var entity T
	if err := e.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (e *Engine[T]) Update(ctx context.Context, entity *T) error {
	return e.db.WithContext(ctx).Save(entity).Error
}

func (e *Engine[T]) Delete(ctx context.Context, id any) error {
	var entity T
	return e.db.WithContext(ctx).Delete(&entity, id).Error
}

func (e *Engine[T]) List(ctx context.Context, opts ...QueryOption) ([]*T, error) {
	var entities []*T
	db := e.applyOptions(e.db.WithContext(ctx), opts)
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (e *Engine[T]) Count(ctx context.Context, opts ...QueryOption) (int64, error) {
	var count int64
	db := e.applyOptions(e.db.WithContext(ctx), opts)
	if err := db.Model(new(T)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (e *Engine[T]) Query(ctx context.Context, opts ...QueryOption) *gorm.DB {
	return e.applyOptions(e.db.WithContext(ctx), opts)
}

func (e *Engine[T]) applyOptions(db *gorm.DB, opts []QueryOption) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (e *Engine[T]) GetByCondition(ctx context.Context, opts ...QueryOption) (*T, error) {
	var entity T
	db := e.applyOptions(e.db.WithContext(ctx), opts)
	if err := db.First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (e *Engine[T]) BatchCreate(ctx context.Context, entities []*T, batchSize ...int) error {
	if len(entities) == 0 {
		return nil
	}

	db := e.db.WithContext(ctx)

	// 设置批量大小
	if len(batchSize) > 0 && batchSize[0] > 0 {
		db = db.CreateInBatches(entities, batchSize[0])
	} else {
		db = db.Create(entities)
	}

	return db.Error
}

func (e *Engine[T]) BatchUpdate(ctx context.Context, updates map[string]interface{}, opts ...QueryOption) (int64, error) {
	if len(updates) == 0 {
		return 0, nil
	}

	db := e.applyOptions(e.db.WithContext(ctx), opts)
	result := db.Model(new(T)).Updates(updates)
	return result.RowsAffected, result.Error
}

func (e *Engine[T]) UpdateByID(ctx context.Context, id any, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	return e.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(updates).Error
}

func (e *Engine[T]) DeleteByCondition(ctx context.Context, opts ...QueryOption) (int64, error) {
	db := e.applyOptions(e.db.WithContext(ctx), opts)
	result := db.Delete(new(T))
	return result.RowsAffected, result.Error
}
func (e *Engine[T]) Upsert(ctx context.Context, entity *T, conflictColumns []string) error {
	if len(conflictColumns) == 0 {
		return e.Create(ctx, entity)
	}

	// 创建冲突处理子句
	conflictClause := clause.OnConflict{
		Columns: make([]clause.Column, 0, len(conflictColumns)),
		DoUpdates: clause.Assignments(
			e.getUpdateAssignments(entity, conflictColumns),
		),
	}

	return e.db.WithContext(ctx).
		Clauses(conflictClause).
		Create(entity).
		Error
}

func (e *Engine[T]) GetByIDWithPreloads(ctx context.Context, id any, preloads ...string) (*T, error) {
	var entity T
	db := e.db.WithContext(ctx)

	// 应用预加载
	for _, p := range preloads {
		if p != "" {
			db = db.Preload(p)
		}
	}

	if err := db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (e *Engine[T]) NewTransaction(ctx context.Context, fn func(EngineInterface[T]) error) error {
	return e.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		engine := NewEngine[T](tx)
		return fn(engine)
	})
}

// getUpdateAssignments 获取需要更新的字段（排除冲突列和创建时间）
var typeCache sync.Map // 缓存结构体字段信息

// getUpdateAssignments 获取更新字段的 map
func (e *Engine[T]) getUpdateAssignments(entity *T, conflictColumns []string) map[string]interface{} {
	modelType := reflect.TypeOf(*entity)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// 缓存 key
	cacheKey := modelType.PkgPath() + "." + modelType.Name()
	conflictSet := make(map[string]struct{}, len(conflictColumns))
	for _, col := range conflictColumns {
		conflictSet[col] = struct{}{}
	}

	// 从缓存获取字段信息
	type fieldInfo struct {
		ColName     string
		IsConflict  bool
		IsCreatedAt bool
		Index       int
	}

	var fields []fieldInfo
	if v, ok := typeCache.Load(cacheKey); ok {
		fields = v.([]fieldInfo)
	} else {
		fields = make([]fieldInfo, 0, modelType.NumField())
		for i := 0; i < modelType.NumField(); i++ {
			f := modelType.Field(i)
			dbTag := f.Tag.Get("gorm")
			if dbTag == "" {
				continue
			}
			colName := e.parseGormTag(dbTag)
			if colName == "" {
				continue
			}
			isCreated := strings.Contains(strings.ToLower(f.Name), "createdat") ||
				strings.Contains(strings.ToLower(dbTag), "autocreatetime")

			fields = append(fields, fieldInfo{
				ColName:     colName,
				IsCreatedAt: isCreated,
				Index:       i,
			})
		}
		typeCache.Store(cacheKey, fields)
	}

	// 收集更新字段
	entityValue := reflect.ValueOf(entity).Elem()
	assignmentsMap := make(map[string]interface{}, len(fields))
	for _, f := range fields {
		if _, ok := conflictSet[f.ColName]; ok || f.IsCreatedAt {
			continue
		}
		assignmentsMap[f.ColName] = entityValue.Field(f.Index).Interface()
	}

	return assignmentsMap
}

// parseGormTag 解析GORM标签获取列名
func (e *Engine[T]) parseGormTag(tag string) string {
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return ""
}
