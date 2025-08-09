package rdb

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EngineInterface 泛型引擎接口
type EngineInterface[T any] interface {
	DB() *gorm.DB
	NewTransaction(ctx context.Context, fn func(EngineInterface[T]) error) error
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id any) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id any) error
	List(ctx context.Context, opts ...QueryOption) ([]*T, error)
	Count(ctx context.Context, opts ...QueryOption) (int64, error)
	Query(ctx context.Context, opts ...QueryOption) *gorm.DB
	GetByCondition(ctx context.Context, opts ...QueryOption) (*T, error)
	BatchCreate(ctx context.Context, entities []*T, batchSize ...int) error
	BatchUpdate(ctx context.Context, updates map[string]interface{}, opts ...QueryOption) (int64, error)
	UpdateByID(ctx context.Context, id any, updates map[string]interface{}) error
	DeleteByCondition(ctx context.Context, opts ...QueryOption) (int64, error)
	Upsert(ctx context.Context, entity *T, conflictColumns []string) error
	GetByIDWithPreloads(ctx context.Context, id any, preloads ...string) (*T, error)
}

// QueryOption 查询选项类型
type QueryOption func(*gorm.DB) *gorm.DB

// ===== 查询选项实现 =====

// WithConditions 条件查询选项
func WithConditions(conds ...any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if len(conds) == 0 {
			return db
		}
		return db.Where(conds[0], conds[1:]...)
	}
}

// WithPagination 分页查询选项
func WithPagination(page, size int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if size <= 0 {
			size = 10
		}
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

// WithOrder 排序选项
func WithOrder(order string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if order != "" {
			return db.Order(order)
		}
		return db
	}
}

// WithPreloads 预加载关联选项
func WithPreloads(preloads ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		for _, p := range preloads {
			if p != "" {
				db = db.Preload(p)
			}
		}
		return db
	}
}

// WithSelect 字段选择选项
func WithSelect(fields ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if len(fields) > 0 {
			return db.Select(fields)
		}
		return db
	}
}

// WithPreloadConditions 带条件的预加载
func WithPreloadConditions(preloadConditions map[string]interface{}) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		for assoc, cond := range preloadConditions {
			if assoc == "" {
				continue
			}

			switch v := cond.(type) {
			case string:
				db = db.Preload(assoc, v)
			case []interface{}:
				if len(v) == 0 {
					db = db.Preload(assoc)
				} else {
					condStr, ok := v[0].(string)
					if !ok {
						continue
					}
					args := append([]interface{}{condStr}, v[1:]...)
					db = db.Preload(assoc, args...)
				}
			default:
				// 判断是否是 gorm 的 Associations 常量
				if cond == clause.Associations {
					db = db.Preload(assoc, clause.Associations)
				} else {
					db = db.Preload(assoc, v)
				}
			}
		}
		return db
	}
}

func WithExclude[T any](fields ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if len(fields) == 0 {
			return db
		}

		var model T
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(&model); err != nil {
			return db
		}

		if stmt.Schema == nil {
			return db
		}

		// 收集所有字段
		allFields := make([]string, 0, len(stmt.Schema.Fields))
		for _, field := range stmt.Schema.Fields {
			if field.DBName != "" {
				allFields = append(allFields, field.DBName)
			}
		}

		// 排除指定字段
		excludeMap := make(map[string]struct{}, len(fields))
		for _, f := range fields {
			excludeMap[f] = struct{}{}
		}

		selectedFields := make([]string, 0, len(allFields))
		for _, f := range allFields {
			if _, excluded := excludeMap[f]; !excluded {
				selectedFields = append(selectedFields, f)
			}
		}

		if len(selectedFields) > 0 {
			return db.Select(selectedFields)
		}
		return db
	}
}

// WithOrderDesc 按字段降序排序
func WithOrderDesc(field string) QueryOption {
	return WithOrder(field + " DESC")
}

// WithOrderAsc 按字段升序排序
func WithOrderAsc(field string) QueryOption {
	return WithOrder(field + " ASC")
}

// WithSQL 自定义SQL条件
func WithSQL(sql string, values ...interface{}) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if sql != "" {
			return db.Where(sql, values...)
		}
		return db
	}
}

// WithLocking 锁选项（悲观锁）
func WithLocking(lockType string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if lockType != "" {
			return db.Clauses(clause.Locking{Strength: lockType})
		}
		return db
	}
}
