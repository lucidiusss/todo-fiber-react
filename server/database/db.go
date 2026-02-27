package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lucidiusss/todo-fiber-react/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	// Load database configuration from environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// Construct DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	// Connect to database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("❌ Failed to get database instance:", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(getEnvAsInt("DB_MAX_IDLE_CONNS", 10))
	sqlDB.SetMaxOpenConns(getEnvAsInt("DB_MAX_OPEN_CONNS", 100))
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("❌ Failed to ping database:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL successfully!")

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}

	fmt.Println("✅ Database migration completed!")

	// Log connection pool stats
	stats := sqlDB.Stats()
	fmt.Printf("📊 Connection Pool: MaxOpen=%d, Open=%d, InUse=%d, Idle=%d\n",
		stats.MaxOpenConnections, stats.OpenConnections, stats.InUse, stats.Idle)
}

// Helper function to get integer environment variables
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
