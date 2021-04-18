package deploy

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mazxaxz/remitly-cli/pkg/remitly"
)

type Code int

const (
	CodeSuccess Code = iota + 1
	CodeError
	CodeTimeout
	CodeUnhealthy
)

func (c *cmdContext) orchestrate(ctx context.Context, rc remitly.Clienter, lbName string, replicas int, result chan Code) {
	for {
		select {
		case <-ctx.Done():
			result <- CodeTimeout
			return
		default:
			time.Sleep(5 * time.Second)
			ss, err := snapshot(ctx, rc, lbName)
			if err != nil {
				result <- CodeError
				return
			}
			if replicas <= 0 {
				for _, ins := range ss.instances {
					if err := rc.DeleteInstance(ctx, lbName, ins.ID); err != nil {
						f := log.Fields{"name": lbName, "id": ins.ID}
						log.WithContext(ctx).WithFields(f).WithError(err).Warn("could not remove instance, skipping")

						result <- CodeError
						return
					}
				}
				result <- CodeSuccess
				return
			}

			finished := true
			for _, instance := range ss.instances {
				if instance.Version != c.revision {
					finished = false
				}
			}
			if finished {
				result <- CodeSuccess
				return
			}

			original := make([]string, 0)
			for _, instance := range ss.instances {
				if instance.Version != c.revision {
					original = append(original, instance.ID)
					continue
				}

				switch instance.Status {
				case remitly.StateProvisioning:
					continue
				case remitly.StateUnhealthy:
					result <- CodeUnhealthy
					return
				case remitly.StateHealthy:
					var ID string
					if len(original) == 0 {
						result <- CodeSuccess
						return
					} else if len(original) == 1 {
						ID, original = original[0], []string{}
					} else {
						ID, original = original[0], original[1:]
					}

					if err := rc.DeleteInstance(ctx, lbName, ID); err != nil {
						result <- CodeError
						return
					}
				}
			}
		}
	}
}
