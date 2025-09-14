package config

import (
	"errors"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ErrNoDSN = errors.New("MYSQL_DSN is empty")

func DBConnection() (*gorm.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")

	if dsn == "" {
		return nil, ErrNoDSN
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger: logger.New(
			log.New(os.Stdout, "", log.LstdFlags),
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				// LogLevel:      logger.Info,
				LogLevel: logger.Warn,
			},
		),
	})

	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("DB Connected")
	return db, nil
}
