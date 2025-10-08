package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the database connection pool
type DB struct {
	*gorm.DB
}

// NewDB creates a new database connection
func NewDB(databaseURL string, logLevel logger.LogLevel) (*DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(postgres.Open(databaseURL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established")

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate runs automatic migrations for all models
func (db *DB) AutoMigrate() error {
	return db.DB.AutoMigrate(
		&User{},
		&Product{},
		&CartItem{},
		&Order{},
		&OrderItem{},
	)
}

// Ping checks if the database connection is alive
func (db *DB) Ping(ctx context.Context) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// WithTransaction executes a function within a database transaction
func (db *DB) WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
