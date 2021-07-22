package domain

import (
	"context"
)

type Probe interface {
	Log(ctx context.Context, value string)
	MonitorDataSegment(ctx context.Context, collection, operation, query string) Segment
}

type Segment interface {
	End()
}
