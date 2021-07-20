package domain

import (
	"context"
)

type Instrumentation interface {
	Log(ctx context.Context, value string)
	// Measure(ctx context.Context, value string)
	// Track(ctx context.Context, value string)
	// NewSegment(ctx context.Context, collection, operation, query string)
}
