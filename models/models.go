package models

import (
	"context"
	"fmt"
	"os"
	"encoding/json"

	"github.com/jackc/pgx/v4"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	connString = "postgress://webapp:webapp@localhost:5432/webapp?sslmode=disable"
)

var (
	db *pgx.Conn
)

func init() {
		db = Connect(connString)
	defer err := db.Close(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not close database connection")
	}

	m, err := migrate.New(
		"file://models/migrations",
		connString)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

func Connect (cstring string) (conn *pgx.Conn) {
	var err error
	config, err := ParseConfig(cstring)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse connection string: %v\n", err)
	}
	conn, err := pgx.Connect(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish database connection: %v\n", err)
		os.Exit(1)
	}
	return &conn
}


// ContactFavorites is a field that contains a contact's favorites
type ContactFavorites struct {
    Colors []string `json:"colors"`
}

// Contact represents a Contact model in the database    
type Contact struct {
	ID                   int
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`

    Name, Address, Phone string

    FavoritesJSON types.JSONText    `db:"favorites"`
    Favorites     *ContactFavorites `db:"-"`
}

type Contacts struct {
	Cs []*Contact{}
}

func (cs *Contacts) Get() error {
    err := db.Query(context.Background(), "select * from contacts").Scan(&cs)
    if err != nil {
        return nil, errors.Wrap(err, "Unable to fetch contacts")
    }

    for _, c := range *cs {
        err := json.Unmarshal(*c.FavoritesJSON, *c.Favorites)

        if err != nil {
            return nil, errors.Wrap(err, "Unable to parse JSON favorites")
        }
    }

    return nil
}