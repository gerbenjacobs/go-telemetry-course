package handler

import (
	"encoding/json"
	"net/http"

	app "github.com/gerbenjacobs/go-telemetry-course"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (h *Handler) schoolTick(w http.ResponseWriter, r *http.Request) {
	schoolCtx, span := h.Tracer.Start(r.Context(), "schoolTick",
		trace.WithAttributes(attribute.String("school.name", h.GoatSchool.SchoolName)))
	defer span.End()

	// get the current time
	now := h.GoatSchool.TickHour
	span.SetAttributes(attribute.Int("school.hour", now))

	h.GoatSchool.Tick()

	// get the class
	class := h.GoatSchool.GetCurrentClass(now)
	if class == nil {
		span.SetStatus(codes.Error, "could not find school class")
		http.Error(w, "No class found for the current hour", http.StatusNotFound)
		return
	}

	// get the students in the class
	students := h.GoatSchool.GetStudentsInClass(class.Name)

	// have students attend class
	h.ClassSvc.AttendClass(schoolCtx, *class, students)

	// output struct
	type GoatClass struct {
		Class    app.GoatClass
		Students []app.GoatStudent
	}

	// write the response
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ") // pretty-print json
	if err := enc.Encode(GoatClass{
		Class:    *class,
		Students: students,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
