package models

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
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

	dbmap.AddTableWithName(Message{}, "messages").SetKeys(true, "Id")
	dbmap.AddTableWithName(Nonce{}, "nonces").SetKeys(true, "Id")
	dbmap.AddTableWithName(Room{}, "rooms").SetKeys(true, "Id")
	dbmap.AddTableWithName(Team{}, "teams").SetKeys(true, "Id")
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(RoomMembership{}, "room_memberships").SetKeys(true, "Id")
	dbmap.AddTableWithName(MessageHeart{}, "message_hearts").setKeys(false, "UserId", "MessageId")

	if os.Getenv("DEBUG") != "" {
		dbmap.TraceOn("[gorp]", log.New(os.Stdout, "[DATABASE] ", log.Lmicroseconds))
	}
	return dbmap
}
