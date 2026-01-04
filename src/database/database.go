package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() error {
	host := getEnv("MYSQL_HOST", "db")
	port := getEnv("MYSQL_PORT", "3306")
	user := getEnv("MYSQL_USER", "quizzler")
	password := getEnv("MYSQL_PASSWORD", "quizzler")
	dbname := getEnv("MYSQL_DATABASE", "quizzler")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, password, host, port, dbname)

	var err error
	for i := 0; i < 30; i++ {
		DB, err = sql.Open("mysql", dsn)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				slog.Info("Connected to database")
				return nil
			}
		}
		slog.Debug("Waiting for database", "attempt", i+1, "max_attempts", 30)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to connect to database after 30 attempts: %v", err)
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
