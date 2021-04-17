package deploy

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/mazxaxz/remitly-cli/pkg/remitly"
	mockRemitly "github.com/mazxaxz/remitly-cli/pkg/remitly/mocks"
)

func TestSnapshot(t *testing.T) {
	t.Run("should create new load balancer and return empty snapshot", func(t *testing.T) {
		// arrange
		const loadBalancerName = "lb_1"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		// expected calls
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(nil, remitly.ErrNotFound)
		mockRemitlyClient.EXPECT().CreateLoadBalancer(gomock.Any(), loadBalancerName).Return(remitly.LoadBalancer{}, nil)

		// act
		result, err := snapshot(context.Background(), mockRemitlyClient, loadBalancerName)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, loadBalancerName, result.loadBalancer)
		assert.Len(t, result.instances, 0)
	})

	t.Run("should return error when load balancer creation fails", func(t *testing.T) {
		// arrange
		const loadBalancerName = "lb_1"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		// expected calls
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(nil, remitly.ErrNotFound)
		mockRemitlyClient.EXPECT().CreateLoadBalancer(gomock.Any(), loadBalancerName).Return(remitly.LoadBalancer{}, remitly.ErrUnknown)

		// act
		result, err := snapshot(context.Background(), mockRemitlyClient, loadBalancerName)

		// assert
		assert.Error(t, err, remitly.ErrUnknown)
		assert.Equal(t, Snapshot{}, result)
	})

	t.Run("should return error when get instances fails", func(t *testing.T) {
		// arrange
		const loadBalancerName = "lb_1"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		// expected calls
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(nil, remitly.ErrUnknown)

		// act
		result, err := snapshot(context.Background(), mockRemitlyClient, loadBalancerName)

		// assert
		assert.Error(t, err, remitly.ErrUnknown)
		assert.Equal(t, Snapshot{}, result)
	})

	t.Run("should return snapshot", func(t *testing.T) {
		// arrange
		const loadBalancerName = "lb_1"
		instances := []remitly.Instance{
			{ID: "i_1", Status: remitly.StateHealthy, Version: "1"},
			{ID: "i_2", Status: remitly.StateHealthy, Version: "1"},
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		// expected calls
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(instances, nil)

		// act
		result, err := snapshot(context.Background(), mockRemitlyClient, loadBalancerName)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, loadBalancerName, result.loadBalancer)
		assert.Equal(t, instances, result.instances)
	})
}
