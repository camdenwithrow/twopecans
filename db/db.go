package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type SqliteConfig struct {
	BaseUrl string
	Token   string
}

func NewSqliteDB(config SqliteConfig) (*sql.DB, error) {
	url := fmt.Sprintf("%s?authToken=%s", config.BaseUrl, config.Token)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", config.BaseUrl, err)
		log.Fatal(err)
	}

	return db, nil
}
