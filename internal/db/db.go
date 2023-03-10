package db

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/raphaelmb/go-grpc-rockets/internal/rocket"
	uuid "github.com/satori/go.uuid"
)

type Store struct {
	db *sqlx.DB
}

func New() (Store, error) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbTable := os.Getenv("DB_TABLE")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbHost, dbPort, dbUsername, dbTable, dbPassword, dbSSLMode,
	)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return Store{}, err
	}
	return Store{
		db: db,
	}, nil
}

// GetRocketByID - retrieves a rocket from the database by id
func (s Store) GetRocketByID(id string) (rocket.Rocket, error) {
	var rkt rocket.Rocket
	row := s.db.QueryRow(`select id, name, type from rockets where id = $1`, id)
	err := row.Scan(&rkt.ID, &rkt.Name, &rkt.Type)
	if err != nil {
		log.Print(err.Error())
		return rocket.Rocket{}, nil
	}
	return rkt, nil
}

// InsertRocket - inserts rocket into the rockets table
func (s Store) InsertRocket(rkt rocket.Rocket) (rocket.Rocket, error) {
	_, err := s.db.NamedQuery(`insert into rockets (id, name, type) values (:id, :name, :type)`, rkt)
	if err != nil {
		return rocket.Rocket{}, errors.New("failed to insert into database")
	}

	return rocket.Rocket{
		ID:   rkt.ID,
		Name: rkt.Name,
		Type: rkt.Type,
	}, nil
}

// DeleteRocket - attempts to delete a rocket from the database return err if error
func (s Store) DeleteRocket(id string) error {
	uid, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`delete from rockets where id = $1`, uid)
	if err != nil {
		return err
	}

	return nil
}
