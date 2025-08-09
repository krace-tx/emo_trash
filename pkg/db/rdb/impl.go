// Package rdb 提供通用的关系型数据库访问层实现
package rdb

import (
	"context"
	"reflect"
	"strings"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Engine 泛型CRUD接口的GORM实现，封装通用数据库操作
// 支持基于GORM的实体增删改查、事务管理、关联查询等功能
type Engine[T any] struct {
	db *gorm.DB // 底层GORM数据库连接实例
}

// NewEngine 创建新的泛型数据库操作实例
// db: GORM数据库连接实例（通常为全局共享连接）
// 返回实现EngineInterface的实例
func NewEngine[T any](db *gorm.DB) EngineInterface[T] {
	return &Engine[T]{db: db}
}

// DB 返回原始的gorm.DB实例，用于执行自定义GORM操作
func (e *Engine[T]) DB() *gorm.DB {
	return e.db
}

// Create 插入单条记录（基于GORM的Create方法）
func (e *Engine[T]) Create(ctx context.Context, entity *T) error {
	return e.db.WithContext(ctx).Create(entity).Error
}

// GetByID 根据主键ID查询单条记录（基于GORM的First方法）
func (e *Engine[T]) GetByID(ctx context.Context, id any) (*T, error) {
	var entity T
	if err := e.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update 全字段更新单条记录（基于GORM的Save方法，会更新所有字段）
func (e *Engine[T]) Update(ctx context.Context, entity *T) error {
	return e.db.WithContext(ctx).Save(entity).Error
}

// Delete 根据主键ID物理删除单条记录（基于GORM的Delete方法）
func (e *Engine[T]) Delete(ctx context.Context, id any) error {
	var entity T
	return e.db.WithContext(ctx).Delete(&entity, id).Error
}

// List 条件查询多条记录，支持分页、排序、预加载等选项
func (e *Engine[T]) List(ctx context.Context, opts ...QueryOption) ([]*T, error) {
	var entities []*T
	db := e.applyOptions(e.db.WithContext(ctx), opts)
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Count 统计符合条件的记录总数（忽略分页条件）
func (e *Engine[T]) Count(ctx context.Context, opts ...QueryOption) (int64, error) {
	var count int64
	db := e.applyOptions(e.db.WithContext(ctx), opts)
	// 使用Model(new(T))确保即使opts包含查询条件，也只统计当前实体表
	if err := db.Model(new(T)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Query 构建带选项的查询器，返回可链式调用的gorm.DB实例
func (e *Engine[T]) Query(ctx context.Context, opts ...QueryOption) *gorm.DB {
	return e.applyOptions(e.db.WithContext(ctx), opts)
}

// applyOptions 应用查询选项到gorm.DB实例，内部辅助方法
func (e *Engine[T]) applyOptions(db *gorm.DB, opts []QueryOption) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

// GetByCondition 根据条件查询单条记录，返回第一条匹配结果
// 若记录不存在且无其他错误，返回(nil, nil)
func (e *Engine[T]) GetByCondition(ctx context.Context, opts ...QueryOption) (*T, error) {
	var entity T
	db := e.applyOptions(e.db.WithContext(ctx), opts)
	if err := db.First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 明确处理记录不存在场景
		}
		return nil, err
	}
	return &entity, nil
}

// BatchCreate 批量插入记录，支持自定义批次大小
// entities: 待插入实体切片（长度为0时直接返回）
// batchSize: 可选批次大小，如BatchCreate(ctx, entities, 100)表示每100条一批插入
func (e *Engine[T]) BatchCreate(ctx context.Context, entities []*T, batchSize ...int) error {
	if len(entities) == 0 {
		return nil
	}

	db := e.db.WithContext(ctx)

	// 根据传入的batchSize决定插入方式
	if len(batchSize) > 0 && batchSize[0] > 0 {
		db = db.CreateInBatches(entities, batchSize[0])
	} else {
		db = db.Create(entities) // 使用GORM默认批次大小
	}

	return db.Error
}

// BatchUpdate 批量更新符合条件的记录（部分字段更新）
// updates: 更新字段映射（key为数据库字段名，value为更新值）
// opts: 查询选项（需包含条件选项，如WithConditions，否则更新全表！）
// 返回受影响的行数和错误
func (e *Engine[T]) BatchUpdate(ctx context.Context, updates map[string]interface{}, opts ...QueryOption) (int64, error) {
	if len(updates) == 0 {
		return 0, nil // 无更新字段时直接返回
	}

	db := e.applyOptions(e.db.WithContext(ctx), opts)
	result := db.Model(new(T)).Updates(updates) // 使用Model确保更新目标表正确
	return result.RowsAffected, result.Error
}

// UpdateByID 根据主键ID更新指定字段（部分字段更新）
// id: 主键值
// updates: 更新字段映射（key为数据库字段名）
func (e *Engine[T]) UpdateByID(ctx context.Context, id any, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil // 无更新字段时直接返回
	}

	return e.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(updates).Error
}

// DeleteByCondition 根据条件批量删除记录（物理删除）
// opts: 查询选项（需包含条件选项，否则删除全表！）
// 返回受影响的行数和错误
func (e *Engine[T]) DeleteByCondition(ctx context.Context, opts ...QueryOption) (int64, error) {
	db := e.applyOptions(e.db.WithContext(ctx), opts)
	result := db.Delete(new(T)) // 使用Model(new(T))确保删除目标表正确
	return result.RowsAffected, result.Error
}

// Upsert 插入或更新记录（基于唯一约束冲突处理）
// entity: 待插入/更新的实体
// conflictColumns: 唯一约束字段列表（发生冲突时更新非冲突字段）
// 实现逻辑：INSERT ... ON CONFLICT (conflictColumns) DO UPDATE SET ...
func (e *Engine[T]) Upsert(ctx context.Context, entity *T, conflictColumns []string) error {
	if len(conflictColumns) == 0 {
		return e.Create(ctx, entity) // 无冲突字段时直接插入
	}

	// 构建冲突处理子句：冲突时更新非冲突字段和非创建时间字段
	conflictClause := clause.OnConflict{
		Columns: make([]clause.Column, 0, len(conflictColumns)),
		DoUpdates: clause.Assignments(
			e.getUpdateAssignments(entity, conflictColumns), // 获取需更新的字段
		),
	}
	// 添加冲突字段
	for _, col := range conflictColumns {
		conflictClause.Columns = append(conflictClause.Columns, clause.Column{Name: col})
	}

	return e.db.WithContext(ctx).
		Clauses(conflictClause).
		Create(entity).
		Error
}

// GetByIDWithPreloads 根据主键ID查询记录并预加载指定关联
// id: 主键值
// preloads: 关联字段列表（如"Orders", "Profile"）
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

// Transaction 开启事务并执行回调函数
// fn: 事务内执行的业务逻辑（参数为事务内的EngineInterface实例）
// 实现逻辑：基于GORM的Transaction方法，自动处理事务提交/回滚
func (e *Engine[T]) Transaction(ctx context.Context, fn func(EngineInterface[T]) error) error {
	return e.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建事务内的Engine实例（使用事务连接tx）
		engine := NewEngine[T](tx)
		return fn(engine) // 执行事务逻辑
	})
}

// getUpdateAssignments 获取Upsert时需更新的字段（排除冲突字段和创建时间）
// 内部使用反射解析实体字段，并缓存结构体信息提升性能
var typeCache sync.Map // 缓存结构体字段信息（key: 结构体类型路径，value: []fieldInfo）

// fieldInfo 结构体字段元信息（用于缓存）
type fieldInfo struct {
	ColName     string // 数据库字段名
	IsConflict  bool   // 是否为冲突字段
	IsCreatedAt bool   // 是否为创建时间字段（需排除更新）
	Index       int    // 字段在结构体中的索引
}

// getUpdateAssignments 解析实体字段，生成Upsert冲突时的更新字段映射
func (e *Engine[T]) getUpdateAssignments(entity *T, conflictColumns []string) map[string]interface{} {
	modelType := reflect.TypeOf(*entity)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// 构建缓存key（结构体包路径+名称）
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
