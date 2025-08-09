package rdb

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EngineInterface 泛型数据库操作接口，定义通用CRUD方法
// 支持单表操作、条件查询、预加载关联、事务等核心功能
// T 为实体类型参数，需对应数据库表结构
type EngineInterface[T any] interface {
	// DB 返回原始gorm.DB实例，用于自定义查询
	DB() *gorm.DB

	// Transaction 开启事务并执行回调函数
	Transaction(ctx context.Context, fn func(EngineInterface[T]) error) error

	// Create 插入单条记录（基于GORM的Create方法）
	Create(ctx context.Context, entity *T) error

	// GetByID 根据主键ID查询单条记录
	// ctx: 上下文
	// id: 主键值（支持int/string等类型，需与实体主键类型匹配）
	// 返回查询到的实体指针和错误（如记录不存在返回gorm.ErrRecordNotFound）
	GetByID(ctx context.Context, id any) (*T, error)

	// Update 更新单条记录（全字段更新，基于主键）
	// ctx: 上下文
	// entity: 待更新的实体指针（需包含主键值）
	// 返回错误信息（如记录不存在、无权限等）
	Update(ctx context.Context, entity *T) error

	// Delete 根据主键ID删除单条记录（物理删除）
	// ctx: 上下文
	// id: 主键值
	// 返回错误信息
	Delete(ctx context.Context, id any) error

	// List 条件查询多条记录
	// ctx: 上下文
	// opts: 查询选项（如条件、分页、排序、预加载等）
	// 返回实体切片和错误信息
	List(ctx context.Context, opts ...QueryOption) ([]*T, error)

	// Count 统计符合条件的记录总数
	// ctx: 上下文
	// opts: 查询选项（用于指定统计条件）
	// 返回记录总数和错误信息
	Count(ctx context.Context, opts ...QueryOption) (int64, error)

	// Query 构建带选项的查询器，用于自定义复杂查询
	// ctx: 上下文
	// opts: 查询选项
	// 返回构建后的gorm.DB实例，可继续链式调用GORM方法
	Query(ctx context.Context, opts ...QueryOption) *gorm.DB

	// GetByCondition 根据条件查询单条记录（返回第一条匹配记录）
	// ctx: 上下文
	// opts: 查询选项（需包含条件类选项，如WithConditions）
	// 返回实体指针和错误（如记录不存在返回nil, nil）
	GetByCondition(ctx context.Context, opts ...QueryOption) (*T, error)

	// BatchCreate 批量插入记录
	// ctx: 上下文
	// entities: 待插入的实体切片指针
	// batchSize: 可选参数，指定每批插入数量（默认使用GORM默认值）
	// 返回错误信息
	BatchCreate(ctx context.Context, entities []*T, batchSize ...int) error

	// BatchUpdate 批量更新符合条件的记录
	// ctx: 上下文
	// updates: 更新字段映射（key: 字段名, value: 字段值）
	// opts: 查询选项（用于指定更新条件）
	// 返回受影响的行数和错误信息
	BatchUpdate(ctx context.Context, updates map[string]interface{}, opts ...QueryOption) (int64, error)

	// UpdateByID 根据主键ID更新指定字段（部分更新）
	// ctx: 上下文
	// id: 主键值
	// updates: 更新字段映射
	// 返回错误信息
	UpdateByID(ctx context.Context, id any, updates map[string]interface{}) error

	// DeleteByCondition 批量删除符合条件的记录
	// ctx: 上下文
	// opts: 查询选项（用于指定删除条件）
	// 返回受影响的行数和错误信息
	DeleteByCondition(ctx context.Context, opts ...QueryOption) (int64, error)

	// Upsert 插入或更新记录（基于唯一约束冲突处理）
	// ctx: 上下文
	// entity: 待插入/更新的实体指针
	// conflictColumns: 唯一约束字段列表（发生冲突时更新非冲突字段）
	// 返回错误信息
	Upsert(ctx context.Context, entity *T, conflictColumns []string) error

	// GetByIDWithPreloads 根据主键ID查询记录并预加载关联数据
	// ctx: 上下文
	// id: 主键值
	// preloads: 关联字段列表（如"Orders", "Profile"）
	// 返回包含关联数据的实体指针和错误
	GetByIDWithPreloads(ctx context.Context, id any, preloads ...string) (*T, error)
}

// QueryOption 查询选项类型，基于函数式选项模式，用于动态配置查询参数
// 示例：WithConditions("status = ?", 1), WithPagination(1, 10), WithPreloads("Orders")
type QueryOption func(*gorm.DB) *gorm.DB

// ===== 查询选项实现 =====

// WithConditions 添加查询条件（支持GORM条件语法）
// conds: 条件参数，格式同gorm.Where()，如：
// - WithConditions("name = ?", "test")
// - WithConditions(map[string]interface{}{"status": 1})
// - WithConditions("age > ? AND score < ?", 18, 100)
func WithConditions(conds ...any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if len(conds) == 0 {
			return db
		}
		return db.Where(conds[0], conds[1:]...)
	}
}

// WithPagination 添加分页条件
// page: 页码（从1开始，<=0时默认1）
// size: 每页条数（<=0时默认10）
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

// WithOrder 添加排序条件
// order: 排序字符串，如"id DESC", "created_at ASC, name DESC"
func WithOrder(order string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if order != "" {
			return db.Order(order)
		}
		return db
	}
}

// WithPreloads 预加载关联数据（无关联条件）
// preloads: 关联字段列表，如[]string{"Orders", "Profile"}
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

// WithSelect 指定查询返回的字段
// fields: 字段名列表，如[]string{"id", "name", "created_at"}
func WithSelect(fields ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if len(fields) > 0 {
			return db.Select(fields)
		}
		return db
	}
}

// WithPreloadConditions 带条件的预加载关联数据
// preloadConditions: 关联条件映射，key为关联字段名，value为条件，支持：
// - 字符串条件："status = ?"
// - 带参数切片：[]interface{}{"amount > ?", 100}
// - 函数条件：func(db *gorm.DB) *gorm.DB { return db.Where("status = 1") }
// - map/struct条件：map[string]interface{}{"status": 1}
// 示例：
//
//	WithPreloadConditions(map[string]interface{}{
//	  "Orders": []interface{}{"status = ?", "paid"},
//	  "Profile": func(db *gorm.DB) *gorm.DB { return db.Select("id", "name") },
//	})
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
