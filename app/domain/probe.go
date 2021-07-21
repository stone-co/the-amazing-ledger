package domain

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type Probe interface {
	Log(ctx context.Context, value string)
	// Measure(ctx context.Context, value string)
	// Track(ctx context.Context, value string)
	MonitorDataSegment(ctx context.Context, collection, operation, query string) *newrelic.DatastoreSegment
}
