package deploy

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/mazxaxz/remitly-cli/pkg/remitly"
	mockRemitly "github.com/mazxaxz/remitly-cli/pkg/remitly/mocks"
)

func TestNewCmd(t *testing.T) {
	t.Run("should return command with specific flags initialized", func(t *testing.T) {
		// arrange

		// act
		cmd := NewCmd()

		// assert
		assert.NotNil(t, cmd.Flag("application"))
		assert.NotNil(t, cmd.Flag("revision"))
		assert.NotNil(t, cmd.Flag("replica-count"))
		assert.NotNil(t, cmd.Flag("wait"))
	})
}

func TestDeploy(t *testing.T) {
	t.Run("should deploy n instances of the application", func(t *testing.T) {
		// arrange
		const loadBalancerName = "lb_1"
		const version = "1"
		const replicas = 3
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		// expected calls
		mockRemitlyClient.EXPECT().CreateInstance(gomock.Any(), loadBalancerName, version).Return(remitly.Instance{}, nil).Times(3)

		// act
		err := deploy(context.Background(), mockRemitlyClient, loadBalancerName, version, replicas)

		// assert
		assert.NoError(t, err)
	})

	t.Run("should return error when at least instance creation fails", func(t *testing.T) {
		// arrange
		const loadBalancerName = "lb_1"
		const version = "1"
		const replicas = 3
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		// expected calls
		mockRemitlyClient.EXPECT().CreateInstance(gomock.Any(), loadBalancerName, version).Return(remitly.Instance{}, remitly.ErrForbidden)

		// act
		err := deploy(context.Background(), mockRemitlyClient, loadBalancerName, version, replicas)

		// assert
		assert.Error(t, err, remitly.ErrForbidden)
	})
}

func TestRollback(t *testing.T) {
	t.Run("should do nothing when original and current snapshots are the same", func(t *testing.T) {
		// arrange
		const loadBalancerName = "lb_1"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		give := Snapshot{
			loadBalancer: loadBalancerName,
			instances: []remitly.Instance{
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
			},
		}
		// expected calls
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(give.instances, nil)

		// act
		err := rollback(context.Background(), mockRemitlyClient, give)

		// assert
		assert.NoError(t, err)
	})

	t.Run("should remove newly created instances when snapshot empty", func(t *testing.T) {
		// arrange
		const loadBalancerName = "lb_1"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		give := Snapshot{
			loadBalancer: loadBalancerName,
			instances:    nil,
		}
		fetched := []remitly.Instance{
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
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(fetched, nil)
		mockRemitlyClient.EXPECT().DeleteInstance(gomock.Any(), loadBalancerName, fetched[0].ID).Return(nil)
		mockRemitlyClient.EXPECT().DeleteInstance(gomock.Any(), loadBalancerName, fetched[1].ID).Return(nil)

		// act
		err := rollback(context.Background(), mockRemitlyClient, give)

		// assert
		assert.NoError(t, err)
	})

	t.Run("should restore original state", func(t *testing.T) {
		// arrange
		const loadBalancerName = "lb_1"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRemitlyClient := mockRemitly.NewMockClienter(mockCtrl)

		give := Snapshot{
			loadBalancer: loadBalancerName,
			instances: []remitly.Instance{
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
			},
		}
		fetched := []remitly.Instance{
			{
				ID:      "ins_1",
				Status:  remitly.StateHealthy,
				Version: "1",
			},
			{
				ID:      "ins_1_2",
				Status:  remitly.StateHealthy,
				Version: "2",
			},
			{
				ID:      "ins_2_2",
				Status:  remitly.StateUnhealthy,
				Version: "2",
			},
		}
		// expected calls
		mockRemitlyClient.EXPECT().GetInstances(gomock.Any(), loadBalancerName).Return(fetched, nil)
		mockRemitlyClient.EXPECT().CreateInstance(gomock.Any(), loadBalancerName, give.instances[1].Version).Return(give.instances[1], nil)
		mockRemitlyClient.EXPECT().DeleteInstance(gomock.Any(), loadBalancerName, fetched[1].ID).Return(nil)
		mockRemitlyClient.EXPECT().DeleteInstance(gomock.Any(), loadBalancerName, fetched[2].ID).Return(nil)

		// act
		err := rollback(context.Background(), mockRemitlyClient, give)

		// assert
		assert.NoError(t, err)
	})
}
