package handler

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func (h *Handler) health(w http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprint(w, "OK!")

	counter, err := h.Meter.Float64Counter(
		"myapp_frontpage",
		metric.WithDescription("a simple counter for our frontpage"),
	)
	if err != nil {
		log.Fatal(err)
	}
	counter.Add(req.Context(), 1)
}

func (h *Handler) randomTime(w http.ResponseWriter, r *http.Request) {
	// create a random duration between 100 and 7500 ms
	randomDuration := time.Duration(100+rand.Intn(7500)) * time.Millisecond

	// start a custom trace span
	_, span := h.Tracer.Start(r.Context(), "inside_randomtime", trace.WithAttributes(attribute.String("duration", randomDuration.String())))

	// Write output and flush instantly
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, _ = w.Write([]byte("Random time for sleep: " + randomDuration.String()))
	w.(http.Flusher).Flush()

	// end the span after all the header stuff and writing
	// but before the sleep time..
	span.End()

	// Sleep for the duration amount
	time.Sleep(randomDuration)
}
