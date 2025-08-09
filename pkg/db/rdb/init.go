package rdb

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	M *gorm.DB // mysql
	P *gorm.DB // postgres
	T *gorm.DB // tinydb
)

type DBType string

const (
	DBTypeMySQL    DBType = "mysql"
	DBTypePostgres DBType = "postgres"
	DBTypeTinyDB   DBType = "tinydb"
)

// 通用数据库配置
type DBConfig struct {
	DSN          string          `json:"dsn"`
	MaxOpenConns int             `json:"max_open_conns,optional"`
	MaxIdleConns int             `json:"max_idle_conns,optional"`
	ConnMaxLife  time.Duration   `json:"conn_max_life,optional"`
	LogLevel     logger.LogLevel `json:"log_level,optional"`
	TablePrefix  string          `json:"table_prefix,optional"`
}

// 初始化MySQL数据库
func InitMySQL(cfg DBConfig) error {
	return initDB(mysql.Open(cfg.DSN), cfg, DBTypeMySQL)
}

// 初始化PostgreSQL数据库
func InitPostgres(cfg DBConfig) error {
	return initDB(postgres.Open(cfg.DSN), cfg, DBTypePostgres)
}

func InitTinyDB(cfg DBConfig) error {
	return initDB(sqlite.Open(cfg.DSN), cfg, DBTypeTinyDB)
}

// 通用数据库初始化逻辑
func initDB(dialector gorm.Dialector, cfg DBConfig, dbType DBType) error {
	// 设置默认值
	setDefaults(&cfg)

	// 创建GORM配置
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,
			SingularTable: true,
		},
	}

	// 建立数据库连接
	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return fmt.Errorf("%s connect failed: %w", dbType, err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get %s sqlDB failed: %w", dbType, err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLife)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("%s ping failed: %w", dbType, err)
	}

	// 设置全局实例
	switch dbType {
	case DBTypeMySQL:
		M = db
	case DBTypePostgres:
		P = db
	case DBTypeTinyDB:
		T = db
	}

	fmt.Printf("%s initialized successfully！\n", dbType)
	return nil
}

// 关闭所有数据库连接
func Close() (err error) {
	if M != nil {
		if sqlDB, e := M.DB(); e == nil {
			if e := sqlDB.Close(); e != nil {
				err = fmt.Errorf("mysql close failed: %w", e)
			}
		}
		M = nil
	}

	if P != nil {
		if sqlDB, e := P.DB(); e == nil {
			if e := sqlDB.Close(); e != nil {
				if err != nil {
					err = fmt.Errorf("%v; postgres close failed: %w", err, e)
				} else {
					err = fmt.Errorf("postgres close failed: %w", e)
				}
			}
		}
		P = nil
	}

	return
}

// 设置配置默认值
func setDefaults(cfg *DBConfig) {
	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = 100
	}
	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = 20
	}
	if cfg.ConnMaxLife <= 0 {
		cfg.ConnMaxLife = 30 * time.Minute
	}
	if cfg.LogLevel == 0 {
		cfg.LogLevel = logger.Info
	}
}
