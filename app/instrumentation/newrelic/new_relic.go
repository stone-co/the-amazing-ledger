package newrelic

import (
	"context"
	"fmt"

	"github.com/newrelic/go-agent/v3/integrations/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

func NewApp(appName, licenseKey string, log *logrus.Entry) (*newrelic.Application, error) {
	if appName == "" || licenseKey == "" {
		log.Warnf("empty app name or license key for new relic application: falling back to empty tracer")

		nrApp, err := newrelic.NewApplication(newrelic.ConfigEnabled(false))
		if err != nil {
			return nil, fmt.Errorf("failed to create new relic application: %w", err)
		}

		return nrApp, nil
	}

	log.WithField("appName", appName).Info("starting new relic tracing")

	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		nrlogrus.ConfigLogger(log.Logger),
		func(cfg *newrelic.Config) {
			cfg.ErrorCollector.RecordPanics = true
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new relic application: %w", err)
	}

	return nrApp, nil
}

func NewDatastoreSegment(ctx context.Context, collection, operation, query string) *newrelic.DatastoreSegment {
	txn := newrelic.FromContext(ctx)
	seg := &newrelic.DatastoreSegment{
		Product:            newrelic.DatastorePostgres,
		Collection:         collection,
		Operation:          operation,
		ParameterizedQuery: query,
	}
	seg.StartTime = txn.StartSegmentNow()
	return seg
}
