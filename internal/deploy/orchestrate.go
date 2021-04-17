package deploy

import (
	"context"
	"time"

	"github.com/mazxaxz/remitly-cli/pkg/remitly"
)

type Code int

const (
	CodeSuccess Code = iota + 1
	CodeError
	CodeTimeout
	CodeUnhealthy
)

func (c *cmdContext) orchestrate(ctx context.Context, rc remitly.Clienter, lbName string, result chan Code) {
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

			previous := make([]string, 0)
			for _, instance := range ss.instances {
				if instance.Version != c.revision {
					previous = append(previous, instance.ID)
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
					ID, previous = previous[0], previous[1:]
					if err := rc.DeleteInstance(ctx, lbName, ID); err != nil {
						result <- CodeError
						return
					}
				}
			}
		}
	}
}