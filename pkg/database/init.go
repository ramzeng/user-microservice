package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var database *gorm.DB

func Initialize(config Config) error {
	var err error

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	if database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return err
	}

	sqlDB, err := database.DB()

	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(config.MaxOpenConnections)
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(config.ConnectionMaxIdleSeconds))
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.ConnectionMaxLifetimeSeconds))

	return nil
}
