package deploy

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mazxaxz/remitly-cli/pkg/remitly"
)

type Snapshot struct {
	loadBalancer string
	instances    []remitly.Instance
}

func snapshot(ctx context.Context, rc remitly.Clienter, lb string) (Snapshot, error) {
	instances, err := rc.GetInstances(ctx, lb)
	if err != nil {
		switch err {
		case remitly.ErrNotFound:
			log.WithContext(ctx).WithField("name", lb).Info("load balancer not found, creating right now...")
			if _, err := rc.CreateLoadBalancer(ctx, lb); err != nil {
				log.WithContext(ctx).WithField("name", lb).WithError(err).Error("could not create load balancer")
				return Snapshot{}, err
			}
			log.WithContext(ctx).WithField("name", lb).Info("load balancer successfully created")
		default:
			log.WithContext(ctx).WithField("name", lb).WithError(err).Error("could not get load balancer instances")
			return Snapshot{}, err
		}
	}

	s := Snapshot{loadBalancer: lb, instances: instances}
	return s, nil
}
