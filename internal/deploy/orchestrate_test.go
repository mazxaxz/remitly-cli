package deploy

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/mazxaxz/remitly-cli/pkg/remitly"
	mockRemitly "github.com/mazxaxz/remitly-cli/pkg/remitly/mocks"
)

func TestOrchestrate(t *testing.T) {
	t.Run("should return timeout when context done", func(t *testing.T) {
		t.Parallel()
		// arrange
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// expected calls

		// act
		result := make(chan Code)
		go orchestrate(ctx, mockRemitlyClient, "lb", "1", 1, result)
		code := <-result

		// assert
		assert.Equal(t, CodeTimeout, code)
	})

	t.Run("should return error when snapshoting", func(t *testing.T) {
		t.Parallel()
		// arrange
		const (
			loadBalancerName = "lb_1"
			version          = "1"
			replicas         = 1
		)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// expected calls
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(nil, remitly.ErrUnknown)

		// act
		result := make(chan Code)
		go orchestrate(ctx, mockRemitlyClient, loadBalancerName, version, replicas, result)
		code := <-result

		// assert
		assert.Equal(t, CodeError, code)
	})

	t.Run("should return unhealthy code when at least one of the deployed instances is unhealthy", func(t *testing.T) {
		t.Parallel()
		// arrange
		const (
			loadBalancerName = "lb_1"
			version          = "1"
			replicas         = 1
		)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		instances := []remitly.Instance{
			{
				ID:      "ins_1",
				Status:  remitly.StateUnhealthy,
				Version: version,
			},
			{
				ID:      "ins_2",
				Status:  remitly.StateHealthy,
				Version: version,
			},
		}

		// expected calls
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(instances, nil)

		// act
		result := make(chan Code)
		go orchestrate(ctx, mockRemitlyClient, loadBalancerName, version, replicas, result)
		code := <-result

		// assert
		assert.Equal(t, CodeUnhealthy, code)
	})

	t.Run("should remove all instances when replica count is 0", func(t *testing.T) {
		t.Parallel()
		// arrange
		const (
			loadBalancerName = "lb_1"
			version          = "2"
			replicas         = 0
		)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		instances := []remitly.Instance{
			{
				ID:      "ins_1",
				Status:  remitly.StateHealthy,
				Version: "1",
			},
			{
				ID:      "ins_2",
				Status:  remitly.StateHealthy,
				Version: "1",
			},
		}

		// expected calls
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(instances, nil)
		mockRemitlyClient.EXPECT().DeleteInstance(gomock.Any(), loadBalancerName, instances[0].ID).Return(nil)
		mockRemitlyClient.EXPECT().DeleteInstance(gomock.Any(), loadBalancerName, instances[1].ID).Return(nil)

		// act
		result := make(chan Code)
		go orchestrate(ctx, mockRemitlyClient, loadBalancerName, version, replicas, result)
		code := <-result

		// assert
		assert.Equal(t, CodeSuccess, code)
	})

	t.Run("should immediately succeed deployment when all instances have desired version", func(t *testing.T) {
		t.Parallel()
		// arrange
		const (
			loadBalancerName = "lb_1"
			version          = "1"
			replicas         = 2
		)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		instances := []remitly.Instance{
			{
				ID:      "ins_1",
				Status:  remitly.StateHealthy,
				Version: version,
			},
			{
				ID:      "ins_2",
				Status:  remitly.StateHealthy,
				Version: version,
			},
		}

		// expected calls
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(instances, nil)

		// act
		result := make(chan Code)
		go orchestrate(ctx, mockRemitlyClient, loadBalancerName, version, replicas, result)
		code := <-result

		// assert
		assert.Equal(t, CodeSuccess, code)
	})

	t.Run("should remove all instances of a previous version", func(t *testing.T) {
		t.Parallel()
		// arrange
		const (
			loadBalancerName = "lb_1"
			oldVersion       = "1"
			version          = "2"
			replicas         = 2
		)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		instancesGen1 := []remitly.Instance{
			{
				ID:      "ins_1",
				Status:  remitly.StateHealthy,
				Version: version,
			},
			{
				ID:      "ins_2",
				Status:  remitly.StateHealthy,
				Version: version,
			},
			{
				ID:      "ins_3",
				Status:  remitly.StateHealthy,
				Version: oldVersion,
			},
			{
				ID:      "ins_4",
				Status:  remitly.StateHealthy,
				Version: oldVersion,
			},
		}
		instancesGen2 := []remitly.Instance{
			{
				ID:      "ins_1",
				Status:  remitly.StateHealthy,
				Version: version,
			},
			{
				ID:      "ins_2",
				Status:  remitly.StateHealthy,
				Version: version,
			},
		}

		// expected calls
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(instancesGen1, nil)
		mockRemitlyClient.EXPECT().DeleteInstance(gomock.Any(), loadBalancerName, instancesGen1[2].ID).Return(nil)
		mockRemitlyClient.EXPECT().DeleteInstance(gomock.Any(), loadBalancerName, instancesGen1[3].ID).Return(nil)
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(instancesGen2, nil)

		// act
		result := make(chan Code)
		go orchestrate(ctx, mockRemitlyClient, loadBalancerName, version, replicas, result)
		code := <-result

		// assert
		assert.Equal(t, CodeSuccess, code)
	})
}
