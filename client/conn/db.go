package conn

import (
	"docmate/config"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func ConnectDB() {
	masterConf := config.DB().Master

	slog.Info("connecting to postgres at ", masterConf.Host, ":", masterConf.Port, "...")
	logMode := logger.Silent
	if masterConf.Debug {
		logMode = logger.Info
	}

	masterDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		masterConf.Host, masterConf.Username, masterConf.Password, masterConf.Name, masterConf.Port)

	slog.Info(masterDSN)
	dB, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{
		PrepareStmt: masterConf.PrepareStmt,
		Logger:      logger.Default.LogMode(logMode),
	})
	if err != nil {
		panic(err)
	}

	sqlDb, err := dB.DB()
	if err != nil {
		panic(err)
	}

	if masterConf.MaxIdleConn != 0 {
		sqlDb.SetMaxIdleConns(masterConf.MaxIdleConn)
	}
	if masterConf.MaxOpenConn != 0 {
		sqlDb.SetMaxOpenConns(masterConf.MaxOpenConn)
	}
	if masterConf.MaxLifeTime != 0 {
		sqlDb.SetConnMaxLifetime(masterConf.MaxLifeTime * time.Second)
	}

	db = dB
	slog.Info("postgres connection successful...")
}

func Db() *gorm.DB {
	return db
}
