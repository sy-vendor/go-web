package mysql

import (
	"database/sql"
	"os"
	"strings"
	"sync"
	"time"

	"go-web/ent"

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

func NewMysql(logger *zap.Logger) *ent.Client {
	once.Do(func() {
		var err error
		var d *sql.DB

		if strings.EqualFold(os.Getenv("DATABASE_WITHOUT_LOGGER"), "true") {
			d = sqldb_logger.OpenDriver(os.Getenv("DATABASE_URL"), &mysql.MySQLDriver{}, zapadapter.New(logger))
		} else {
			d, err = sql.Open("mysql", os.Getenv("DATABASE_URL"))
			if err != nil {
				panic(err)
			}
		}

		d.SetConnMaxLifetime(time.Hour)
		d.SetMaxOpenConns(100)
		d.SetConnMaxIdleTime(10)
		db = ent.NewClient(ent.Driver(ent_sql.OpenDB(dialect.MySQL, d)))
	})

	return db
}

func CloseMysql() error {
	return db.Close()
}
