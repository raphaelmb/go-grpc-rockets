// +acceptance

package test

import (
	"context"
	"testing"

	v1 "github.com/raphaelmb/go-grpc-rockets/protos/rocket/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RocketTestSuite struct {
	suite.Suite
}

func (s *RocketTestSuite) TestAddRocket() {
	s.T().Run("adds a new rocket successfully", func(t *testing.T) {
		client := GetClient()
		resp, err := client.AddRocket(context.Background(), &v1.AddRocketRequest{Rocket: &v1.Rocket{
			Id:   "bac7d1eb-9642-4c76-b1ad-6c114ca220f3",
			Name: "V1",
			Type: "Falcon Heavy",
		}})
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "bac7d1eb-9642-4c76-b1ad-6c114ca220f3", resp.Rocket.Id)
	})

	s.T().Run("validates the uuid in the new rocket is a uuid", func(t *testing.T) {
		client := GetClient()
		_, err := client.AddRocket(context.Background(), &v1.AddRocketRequest{Rocket: &v1.Rocket{
			Id:   "not-a-valid-uuid",
			Name: "V1",
			Type: "Falcon Heavy",
		}})
		assert.Error(s.T(), err)

		st := status.Convert(err)
		assert.Equal(s.T(), codes.InvalidArgument, st.Code())
	})
}

func TestRocketService(t *testing.T) {
	suite.Run(t, new(RocketTestSuite))
}
