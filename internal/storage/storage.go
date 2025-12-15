package storage

import (
	"bank_app/internal/config"
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// метод DB для применения миграций
func (d *DataBase) UpMigrations() error {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Printf("Failed to set dialect: %v", err)
		return err
	}

	err := goose.Up(d.DB, d.MigrationsDir)
	if err != nil {
		log.Printf("Failed to upping migrations: %v", err)
		return err
	}

	log.Println("Migrations is upping")
	return nil
}

// соединение с DB и return экземпляра
func NewStorage(cfg *config.Config) (*DataBase, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Failed to open database:", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Failed to ping database:", err)
		return nil, err
	}

	log.Println("Successfully connected to database")

	datebase := DataBase{
		DB:            db,
		MigrationsDir: cfg.Database.MigrationsDir,
	}

	return &datebase, nil
}
