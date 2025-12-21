package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=annygo sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Println("❌ erro ao conectar no banco")
		return nil, err
	}

	log.Println("✅ conexão com o banco realizada com sucesso")
	return db, nil
}
