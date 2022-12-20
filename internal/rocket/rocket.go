//go:generate mockgen -destination=rocket_mocks_test.go -package=rocket github.com/raphaelmb/go-grpc-rockets/internal/rocket Store
package rocket

import "context"

// Contains the definition of the rocket
type Rocket struct {
	ID     string
	Name   string
	Type   string
	Flight int
}

// Interface that the database implementation is expected to follow
type Store interface {
	GetRocketByID(id string) (Rocket, error)
	InsertRocket(rkt Rocket) (Rocket, error)
	DeleteRocket(id string) error
}

// Responsible for updating the rocket inventory
type Service struct {
	Store Store
}

// Returns new instance of rocket service
func New(store Store) Service {
	return Service{
		Store: store,
	}
}

func (s Service) GetRocketByID(ctx context.Context, id string) (Rocket, error) {
	rkt, err := s.Store.GetRocketByID(id)
	if err != nil {
		return Rocket{}, nil
	}
	return rkt, nil
}

func (s Service) InsertRocket(ctx context.Context, rkt Rocket) (Rocket, error) {
	rkt, err := s.Store.InsertRocket(rkt)
	if err != nil {
		return Rocket{}, nil
	}
	return rkt, nil
}

func (s Service) DeleteRocket(id string) error {
	err := s.Store.DeleteRocket(id)
	if err != nil {
		return err
	}
	return nil
}
