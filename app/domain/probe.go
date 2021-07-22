package domain

import (
	"context"
)

type Probe interface {
	Log(ctx context.Context, value string)
	MonitorSegment(ctx context.Context) Segment
	MonitorDataSegment(ctx context.Context, collection, operation, query string) Segment
}

type Segment interface {
	End()
}
