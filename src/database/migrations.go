package database

import (
	"fmt"
	"log/slog"
)

func Migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			admin TINYINT NOT NULL DEFAULT 0,
			active TINYINT NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX email_idx (email),
			INDEX active_idx (active)
		)`,
		`CREATE TABLE IF NOT EXISTS decks (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			active TINYINT NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			INDEX name_idx (name),
			INDEX active_name_idx (active, name)
		)`,
		`CREATE TABLE IF NOT EXISTS cards (
			id INT AUTO_INCREMENT PRIMARY KEY,
			deck_id INT NOT NULL,
			front TEXT NOT NULL,
			back TEXT NOT NULL,
			active TINYINT NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (deck_id) REFERENCES decks(id) ON DELETE CASCADE
		)`,
	}

	for _, migration := range migrations {
		if _, err := DB.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %v", err)
		}
	}

	slog.Info("Migrations completed successfully")
	return nil
}
