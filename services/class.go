package services

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	app "github.com/gerbenjacobs/go-telemetry-course"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ClassService struct {
	tracer trace.Tracer
}

func NewClassService() *ClassService {
	return &ClassService{
		tracer: otel.Tracer("classService"),
	}
}

func (c *ClassService) AttendClass(ctx context.Context, class app.GoatClass, students []app.GoatStudent) {
	classCtx, span := c.tracer.Start(ctx, "attendClass")
	span.SetAttributes(attribute.String("class.name", class.Name))
	defer span.End()

	// Simulate attending class
	for _, student := range students {
		slog.Debug("adding student to class", "student", student.Name, "class", class.Name)
		_, studentSpan := c.tracer.Start(classCtx, "attendStudent",
			trace.WithAttributes(attribute.String("student.name", student.Name)))

		// Simulate some processing
		randomDuration := time.Duration(100+rand.Intn(1000)) * time.Millisecond
		<-time.After(randomDuration)

		studentSpan.End()
	}
}
