package mysql

import (
	"database/sql"
	"sync"

	"go-web/ent"
	"go-web/pkg/config"

	"entgo.io/ent/dialect"
	ent_sql "entgo.io/ent/dialect/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	sqldb_logger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zapadapter"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(NewMysql)

var (
	db   *ent.Client
	once sync.Once
)

func NewMysql(cfg *config.Config, logger *zap.Logger) *ent.Client {
	once.Do(func() {
		var err error
		var d *sql.DB

		// 使用配置中的数据库连接信息
		dbConfig := cfg.Database
		dsn := dbConfig.GetDSN()

		// 根据配置决定是否使用日志记录器
		if cfg.Log.Level == "debug" {
			d = sqldb_logger.OpenDriver(dsn, &mysql.MySQLDriver{}, zapadapter.New(logger))
		} else {
			d, err = sql.Open("mysql", dsn)
			if err != nil {
				logger.Fatal("failed to open database connection",
					zap.Error(err),
					zap.String("dsn", maskDSN(dsn)),
				)
			}
		}

		// 使用配置中的连接池设置
		d.SetConnMaxLifetime(dbConfig.MaxLifetime)
		d.SetMaxOpenConns(dbConfig.MaxOpenConns)
		d.SetMaxIdleConns(dbConfig.MaxIdleConns)

		// 测试数据库连接
		if err := d.Ping(); err != nil {
			logger.Fatal("failed to ping database",
				zap.Error(err),
				zap.String("dsn", maskDSN(dsn)),
			)
		}

		logger.Info("database connection established",
			zap.String("host", dbConfig.Host),
			zap.Int("port", dbConfig.Port),
			zap.String("database", dbConfig.Database),
			zap.Int("max_open_conns", dbConfig.MaxOpenConns),
			zap.Int("max_idle_conns", dbConfig.MaxIdleConns),
			zap.Duration("max_lifetime", dbConfig.MaxLifetime),
		)

		db = ent.NewClient(ent.Driver(ent_sql.OpenDB(dialect.MySQL, d)))
	})

	return db
}

func CloseMysql() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// maskDSN masks sensitive information in DSN
func maskDSN(dsn string) string {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "***"
	}
	cfg.Passwd = "***"
	return cfg.FormatDSN()
}
