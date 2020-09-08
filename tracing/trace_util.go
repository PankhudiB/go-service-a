package tracing

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"net/http"
	"time"
)

func Init(serviceName string, ocAgent string) *ocagent.Exporter {
	oce, _ := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(1*time.Second),
		ocagent.WithAddress(ocAgent),
		ocagent.WithServiceName(serviceName))
	trace.RegisterExporter(oce)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	return oce
}

func WithTracing(app http.Handler) http.Handler {
	return &ochttp.Handler{
		Handler: app,
		GetStartOptions: func(r *http.Request) trace.StartOptions {
			startOptions := trace.StartOptions{}
			return startOptions
		},
	}
}
