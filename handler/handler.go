package handler

import (
	"net/http"

	"github.com/gerbenjacobs/go-telemetry-course/services"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Handler is your dependency container and HTTP handler
type Handler struct {
	mux http.Handler
	Dependencies
}

// Dependencies contains all the dependencies your application and its services require
type Dependencies struct {
	Meter  metric.Meter
	Tracer trace.Tracer

	GoatSchool *services.GoatSchool // using an implementation, instead of interface! >:(
	ClassSvc   *services.ClassService
}

// New creates a new handler given a set of dependencies
func New(dependencies Dependencies) *Handler {
	h := &Handler{
		Dependencies: dependencies,
	}

	// create a new HTTP mux and set routes
	r := http.NewServeMux()

	r.Handle("GET /{$}", otelhttp.WithRouteTag("/", http.HandlerFunc(h.health)))
	r.Handle("GET /randomtime", otelhttp.WithRouteTag("/randomtime", http.HandlerFunc(h.randomTime)))
	r.Handle("GET /school/tick", otelhttp.WithRouteTag("/school/tick", http.HandlerFunc(h.schoolTick)))

	// wrap the mux with OpenTelemetry
	h.mux = otelhttp.NewHandler(r, "myapp")
	return h
}

// ServeHTTP makes sure Handler implements the http.Handler interface
// this keeps the underlying mux private
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
