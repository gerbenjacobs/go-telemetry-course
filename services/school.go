package services

import (
	app "github.com/gerbenjacobs/go-telemetry-course"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type GoatSchool struct {
	SchoolName string
	TickHour   int // represents hours in school from 1 to 8
	Students   []app.GoatStudent
	Classes    []app.GoatClass

	tracer trace.Tracer
}

func NewGoatSchool(name string) *GoatSchool {
	return &GoatSchool{
		SchoolName: name,
		Students:   []app.GoatStudent{},
		Classes:    []app.GoatClass{},

		tracer: otel.Tracer("goatSchool"),
	}
}

func (g *GoatSchool) Tick() {
	// update the hour for next tick
	g.TickHour++
	if g.TickHour > 8 {
		g.TickHour = 1
	}
}

func (g *GoatSchool) AddStudent(name string, age int, classes []string) {
	g.Students = append(g.Students, app.GoatStudent{
		Name:    name,
		Age:     age,
		Classes: classes,
	})
}

func (g *GoatSchool) AddClass(name string, startTime int) {
	g.Classes = append(g.Classes, app.GoatClass{
		Name:      name,
		StartTime: startTime,
	})
}

// GetCurrentClass returns the current class for a given hour
func (g *GoatSchool) GetCurrentClass(hour int) *app.GoatClass {
	for _, class := range g.Classes {
		if class.StartTime == hour {
			return &class
		}
	}
	return nil
}

// GetStudentsInClass returns the students in a given class
func (g *GoatSchool) GetStudentsInClass(className string) []app.GoatStudent {
	var studentsInClass []app.GoatStudent
	for _, student := range g.Students {
		for _, class := range student.Classes {
			if class == className {
				studentsInClass = append(studentsInClass, student)
			}
		}
	}
	return studentsInClass
}
