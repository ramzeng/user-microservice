package database

import "gorm.io/gorm"

func Where(query interface{}, args ...interface{}) *gorm.DB {
	return database.Where(query, args...)
}

func Create(value interface{}) *gorm.DB {
	return database.Create(value)
}

func AutoMigrate(dst ...interface{}) error {
	return database.AutoMigrate(dst...)
}
