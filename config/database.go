package config

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	// Mengambil konfigurasi dari package config
	config := Config

	// Escape password untuk URI PostgreSQL
	encodedPassword := url.QueryEscape(config.Database.Password)

	// Membuat URI untuk koneksi database
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.Database.Username,
		encodedPassword,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	// Membuka koneksi ke database dengan Gorm
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Mendapatkan objek *sql.DB dari Gorm untuk pengaturan koneksi
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Konfigurasi koneksi database
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConnection)
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConnection)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.MaxLifeTimeConnection) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.Database.MaxIdleTime) * time.Second)

	return db, nil
}
