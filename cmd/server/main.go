package main

import (
	"log"

	"github.com/raphaelmb/go-grpc-rockets/internal/db"
	"github.com/raphaelmb/go-grpc-rockets/internal/rocket"
)

func Run() error {
	// responsible for initializing and starting the grpc server
	rocketStore, err := db.New()
	if err != nil {
		return err
	}
	err = rocketStore.Migrate()
	if err != nil {
		log.Println("Failed to run migrations")
		return err
	}
	_ = rocket.New(rocketStore)
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
