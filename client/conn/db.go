package conn

import (
	"fmt"
	"log/slog"
	"time"

	"docmate/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func ConnectDB() {
	dbConf := config.DB().Db

	slog.Info("connecting to postgres",
		"host", dbConf.Host,
		"port", dbConf.Port,
	)
	logMode := logger.Silent
	if dbConf.Debug {
		logMode = logger.Info
	}

	masterDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConf.Host, dbConf.Username, dbConf.Password, dbConf.Name, dbConf.Port)

	slog.Info("postgres dsn", "dsn", masterDSN)
	dB, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{
		PrepareStmt: dbConf.PrepareStmt,
		Logger:      logger.Default.LogMode(logMode),
	})
	if err != nil {
		panic(err)
	}

	sqlDb, err := dB.DB()
	if err != nil {
		panic(err)
	}

	if dbConf.MaxIdleConn != 0 {
		sqlDb.SetMaxIdleConns(dbConf.MaxIdleConn)
	}
	if dbConf.MaxOpenConn != 0 {
		sqlDb.SetMaxOpenConns(dbConf.MaxOpenConn)
	}
	if dbConf.MaxLifeTime != 0 {
		sqlDb.SetConnMaxLifetime(dbConf.MaxLifeTime * time.Second)
	}

	db = dB
	slog.Info("postgres connection successful...")
}

func Db() *gorm.DB {
	return db
}
