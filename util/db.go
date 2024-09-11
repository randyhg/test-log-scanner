package util

import (
	"database/sql"
	"github.com/randyhg/test-log-scanner/config"
	"github.com/randyhg/test-log-scanner/util/mylog"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

func Master() *gorm.DB {
	return db
}

func openDB(dsn string, config *gorm.Config, maxIdleConns, maxOpenConns int) (err error) {
	if config == nil {
		config = &gorm.Config{}
	}

	if config.NamingStrategy == nil {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		}
	}

	if db, err = gorm.Open(mysql.Open(dsn), config); err != nil {
		log.Printf("opens database failed: %v", err.Error())
		return
	}

	if sqlDB, err = db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
	} else {
		log.Println(err)
	}
	return
}

func InitDB() {
	gormConf := &gorm.Config{}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	gormConf.Logger = newLogger
	err := openDB(config.Instance.MySqlUrl, gormConf,
		config.Instance.MySqlMaxIdle, config.Instance.MySqlMaxOpen)
	if err != nil {
		log.Fatal(err)
	}
	mylog.Info("MySQL connection established")
}

func CreateMonthTable(db *gorm.DB, dst schema.Tabler, tableName string) error {
	mig := db.Migrator()
	if !mig.HasTable(tableName) {
		if err := mig.CreateTable(dst); err != nil {
			return err
		}
		if err := mig.RenameTable(dst.TableName(), tableName); err != nil {
			return err
		}
	}
	return nil
}
