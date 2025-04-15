package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gerbenjacobs/go-telemetry-course/handler"
	"github.com/gerbenjacobs/go-telemetry-course/internal"
	"github.com/gerbenjacobs/go-telemetry-course/services"
	"github.com/lmittmann/tint"
	"go.opentelemetry.io/otel"
)

const (
	serviceName    = "github.com/gerbenjacobs/go-telemetry-course"
	serviceVersion = "v0.1.0"
)

func main() {
	// handle shutdown signals
	shutdown := make(chan os.Signal, 3)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// set output logging
	level := slog.LevelDebug
	slog.SetDefault(
		slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{Level: level}),
		),
	)

	// setup OpenTelemetry
	otelShutdown, err := internal.SetupOTelSDK(context.Background(), serviceName, serviceVersion)
	if err != nil {
		slog.Error("Failed to setup OpenTelemetry", "error", err)
		os.Exit(1)
	}

	// set up your goat school!
	goatSchool := services.NewGoatSchool("Kid Valley High")
	goatSchool.AddClass("History", 1)
	goatSchool.AddClass("Biology", 2)
	goatSchool.AddClass("Jumping", 5)
	goatSchool.AddClass("Bleating", 7)

	goatSchool.AddStudent("Billy", 5, []string{"History", "Biology"})
	goatSchool.AddStudent("Vincent van Goat", 6, []string{"History", "Jumping"})
	goatSchool.AddStudent("Scape", 4, []string{"Bleating", "Jumping"})
	goatSchool.AddStudent("Daisy", 5, []string{"Bleating", "Biology"})

	// set up the route handler and server
	app := handler.New(handler.Dependencies{
		Meter:  otel.Meter(serviceName),
		Tracer: otel.Tracer(serviceName),

		GoatSchool: goatSchool,
		ClassSvc:   services.NewClassService(),
	})
	srv := &http.Server{
		Addr:         ":9000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      app,
	}

	// start running the server
	go func() {
		log.Print("Server started on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to listen", "error", err)
			os.Exit(1)
		}
	}()

	// wait for shutdown signals
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	if err := otelShutdown(ctx); err != nil {
		slog.Error("Failed to shutdown OpenTelemetry", "error", err)
		os.Exit(1)
	}
	slog.Info("Server shutdown complete")
}
