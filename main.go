package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Handler struct {
	db *sql.DB
}

func (h *Handler) SingleInsertperRow(link string) error {
	query := "INSERT INTO links (url) VALUES (?)"
	_, err := h.db.Exec(query, link)
	return err
}

func (h *Handler) PreparedStatementsExecute(links []string) error {
	query := "INSERT INTO links (url) VALUES (?)"
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, link := range links {
		if _, err := stmt.Exec(link); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) BatchInsert(links []string) error {
	if len(links) == 0 {
        return nil
    }

	query := "INSERT INTO links (url) VALUES "
	vals := make([]interface{}, 0, len(links))

	for i, link := range links {
		query += "(?)"
		if i < len(links)-1 {
			query += ","
		}
		vals = append(vals, link)
	}

	_, err := h.db.Exec(query, vals...);
	return err

}

func (h *Handler) TransactionInserts(links []string) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO links (url) VALUES (?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, link := range links {
		if _, err := stmt.Exec(link); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (h *Handler) TransactionBatchInserts(links []string) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}

	query := "INSERT INTO links (url) VALUES "
	vals := make([]interface{}, 0, len(links))

	for i, link := range links {
		query += "(?)"
		if i < len(links)-1 {
			query += ","
		}
		vals = append(vals, link)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(vals...); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func main() {
	db, err := connectDb()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
}

func connectDb() (*sql.DB, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "linksdb"
	cfg.ParseTime = true

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to SQL Database!")

	return db, nil
}
