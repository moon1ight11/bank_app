package transactionsservice

import (
	"bank_app/internal/config"
	"bank_app/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
)

var testDB *storage.DataBase

func TestMain(m *testing.M) {
	cfg := config.Config{
		Database: config.DatabaseConfig{
			Host:          "localhost",
			Port:          15432,
			DBName:        "bank_app",
			User:          "fedor",
			Password:      "fedor_pass",
			MigrationsDir: "../../../internal/storage/repos/migrations",
		},
	}

	if err := createTestDatabase(&cfg); err != nil {
		log.Fatalf("Failed to create test database: %v", err)
	}

	db, err := storage.NewStorage(&cfg)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.UpMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	ctx := context.Background()
	query1 := `
			INSERT INTO bank_app.users (id, name, surname, email, phone_number, password, role)
    		VALUES 
				('00000000-0000-0000-0000-000000000001', 
				'admin', 
				'zero', 
				'nulladmin@mail.com', 
				'+70000000000', 
				'$2a$10$rQxxz3Hg/g2IRyIhm8K9/u.q0onCMAve1bGcfbHNeBYbvojhjjKVS', 
				'Admin')
    		ON CONFLICT (id) DO NOTHING`
	db.DB.ExecContext(ctx, query1)

	query2 := `
			INSERT INTO bank_app.accounts (id, user_id, currency)
    		VALUES 
      		  	('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'RUB'),
       		 	('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'EUR'),
        		('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'USD')
    		ON CONFLICT (id) DO NOTHING`
	db.DB.ExecContext(ctx, query2)

	testDB = db

	code := m.Run()

	cleanTestDatabase(db)
	db.DB.Close()
	os.Exit(code)
}

func createTestDatabase(cfg *config.Config) error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.Database.DBName).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Database.DBName))
		if err != nil {
			return err
		}
	}

	return nil
}

func cleanTestDatabase(db *storage.DataBase) {
	ctx := context.Background()
	tables := []string{"bank_app.transactions", "bank_app.accounts", "bank_app.users"}
	for _, table := range tables {
		db.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table))
	}
}
