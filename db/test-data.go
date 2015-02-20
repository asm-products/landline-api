// forego run go run db/test-data.go

package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/asm-products/landline-api/models"
	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "[DATABASE] ", log.Lmicroseconds))

	dbmap.AddTableWithName(models.Team{}, "teams").SetKeys(true, "id")

	team := &models.Team{
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Email:             "jake@ooo.com",
		EncryptedPassword: "s3kr3th4sh",
		OAuthAuthorizeUrl: "http://oooid.com/oauth/authorize",
		OAuthTokenUrl:     "http://oooid.com/oauth/authorize",
		Slug:              "asm-dev",
	}

	err = dbmap.Insert(team)
	if err != nil {
		log.Fatal(err)
	}
}
