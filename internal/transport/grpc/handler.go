package grpc

import (
	"context"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/raphaelmb/go-grpc-rockets/internal/rocket"
	rkt "github.com/raphaelmb/go-grpc-rockets/protos/rocket/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// defines the interface that the concrete implementations has to adhere to
type RocketService interface {
	GetRocketByID(ctx context.Context, id string) (rocket.Rocket, error)
	InsertRocket(ctx context.Context, rkt rocket.Rocket) (rocket.Rocket, error)
	DeleteRocket(ctx context.Context, id string) error
}

// handles grpc requests
type Handler struct {
	RocketService RocketService
	rkt.UnimplementedRocketServiceServer
}

// return a new grpc handler
func New(rktService RocketService) Handler {
	return Handler{
		RocketService: rktService,
	}
}

func (h Handler) Serve() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Print("could not listen on port 50051")
		return err
	}

	grpcServer := grpc.NewServer()
	rkt.RegisterRocketServiceServer(grpcServer, &h)

	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("failed to server %s\n", err)
		return err
	}

	return nil
}

// retrieves a rocket by id and returns the response
func (h Handler) GetRocket(ctx context.Context, req *rkt.GetRocketRequest) (*rkt.GetRocketResponse, error) {
	log.Print("Get Rocket gRPC endpoint hit")

	rocket, err := h.RocketService.GetRocketByID(ctx, req.Id)
	if err != nil {
		log.Print("Failed to retrieve rocket by ID")
		return &rkt.GetRocketResponse{}, nil
	}

	return &rkt.GetRocketResponse{
		Rocket: &rkt.Rocket{
			Id:   rocket.ID,
			Name: rocket.Name,
			Type: rocket.Type,
		},
	}, nil
}

// AddRocket - adds rocket to the database
func (h Handler) AddRocket(ctx context.Context, req *rkt.AddRocketRequest) (*rkt.AddRocketResponse, error) {
	log.Print("Add Rocket gRPC endpoint hit")

	if _, err := uuid.Parse(req.Rocket.Id); err != nil {
		errorStatus := status.Error(codes.InvalidArgument, "uuid is not valid")
		log.Print("given uuid is not valid")
		return &rkt.AddRocketResponse{}, errorStatus
	}

	newRkt, err := h.RocketService.InsertRocket(ctx, rocket.Rocket{
		ID:   req.Rocket.Id,
		Name: req.Rocket.Name,
		Type: req.Rocket.Type,
	})
	if err != nil {
		log.Print("Failed to insert rocket into database")
		return &rkt.AddRocketResponse{}, err
	}
	return &rkt.AddRocketResponse{
		Rocket: &rkt.Rocket{
			Id:   newRkt.ID,
			Name: newRkt.Name,
			Type: newRkt.Type,
		},
	}, nil

}

// DeleteRocket - deletes a rocket from the database
func (h Handler) DeleteRocket(ctx context.Context, req *rkt.DeleteRocketRequest) (*rkt.DeleteRocketResponse, error) {
	log.Print("Delete Rocket gRPC endpoint hit")

	err := h.RocketService.DeleteRocket(ctx, req.Rocket.Id)
	if err != nil {
		log.Print("Failed to delete rocket from database")
		return &rkt.DeleteRocketResponse{}, err
	}
	return &rkt.DeleteRocketResponse{
		Status: "successfully deleted rocket",
	}, nil
}
