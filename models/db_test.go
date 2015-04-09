package models
import (
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
    "gopkg.in/gorp.v1"
)

func newDbTestContext() {
    db, err := sql.Open("sqlite3", ":memory:")
    if (err != nil) {
        log.Fatalln(err)
    }
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

    dbmap.AddTableWithName(Message{}, "messages").SetKeys(false, "Id")
    dbmap.AddTableWithName(Nonce{}, "nonces").SetKeys(false, "Id")
    dbmap.AddTableWithName(Room{}, "rooms").SetKeys(false, "Id")
    dbmap.AddTableWithName(Team{}, "teams").SetKeys(false, "Id")
    dbmap.AddTableWithName(User{}, "users").SetKeys(false, "Id")
    dbmap.AddTableWithName(RoomMembership{}, "room_memberships").SetKeys(false, "Id")

    err = dbmap.CreateTables()
    if (err != nil) {
        log.Fatalln(err)
    }

    Db = dbmap
}