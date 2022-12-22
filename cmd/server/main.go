package main

import (
	"log"

	_ "github.com/gogo/protobuf/protoc-gen-gogo/grpc"
	"github.com/raphaelmb/go-grpc-rockets/internal/db"
	"github.com/raphaelmb/go-grpc-rockets/internal/rocket"
	"github.com/raphaelmb/go-grpc-rockets/internal/transport/grpc"
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

	rktService := rocket.New(rocketStore)
	rktHandler := grpc.New(rktService)

	if err := rktHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
