package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/asm-products/landline-api/models"
	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

var (
	// Db is the global database context
	Db = NewDbContext(os.Getenv("DATABASE_URL"))
)

// NewDbContext initialises a new database context
func NewDbContext(url string) *gorp.DbMap {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(models.Message{}, "messages").SetKeys(false, "Id")
	dbmap.AddTableWithName(models.Room{}, "rooms").SetKeys(false, "Id")
	dbmap.AddTableWithName(models.Team{}, "teams").SetKeys(false, "Id")
	dbmap.AddTableWithName(models.User{}, "users").SetKeys(false, "Id")

	if os.Getenv("DEBUG") != "" {
		dbmap.TraceOn("[gorp]", log.New(os.Stdout, "[DATABASE] ", log.Lmicroseconds))
	}
	return dbmap
}
