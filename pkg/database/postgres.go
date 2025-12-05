package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// InitDB inisialisasi koneksi database dengan GORM
func InitDB(config Config) (*gorm.DB, error) {
	// Format DSN untuk PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
	)

	// Koneksi ke database
	db, err := gorm.Open(postgresdriver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log SQL queries
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✅ Database connected successfully")

	return db, nil
}

// RunMigrations menjalankan database migrations menggunakan golang-migrate
func RunMigrations(config Config) error {
	// Format connection string untuk migrate
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	// Open database connection untuk migrate
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database for migration: %w", err)
	}
	defer sqlDB.Close()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Path ke migration files
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations
	log.Println("🔄 Running database migrations...")
	err = m.Up() // assign ke err yang sudah ada
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("✅ Database schema is up to date")
	} else {
		log.Println("✅ Database migrations completed successfully")
	}

	return nil
}
