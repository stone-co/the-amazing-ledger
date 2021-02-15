package newrelic

import (
	"context"

	"github.com/newrelic/go-agent/v3/integrations/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

func NewRelicApp(appName, licenseKey string, log *logrus.Entry) (*newrelic.Application, error) {
	if appName == "" || licenseKey == "" {
		log.Warnf("empty app name or license key for new relic application: falling back to empty tracer")
		return newrelic.NewApplication(newrelic.ConfigEnabled(false))
	}

	log.WithField("appName", appName).Info("starting new relic tracing")
	return newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		nrlogrus.ConfigLogger(log.Logger),
		func(cfg *newrelic.Config) {
			cfg.ErrorCollector.RecordPanics = true
		},
	)
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
