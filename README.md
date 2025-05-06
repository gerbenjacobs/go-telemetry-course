# Goat Elementary 

## A Go OpenTelemetry course

In this course you'll learn things such as OpenTelemetry, Prometheus, Grafana and
how to implement this in your Go application.

This ~Go Telemetry~ Goat Elementary course has you set up a simple fake Go application
with multiple services to simulate tracing between different services.

You are supposed to implement new **metrics** and **tracing** functionality in the code
and come up with proper dashboards to visualize the data.

## Prerequisites

- a working Go setup
- Docker
- (optionally) `watch` to simulate traffic (`brew install watch`)

## Getting started

Start by running Grafana's OpenTelemetry LGTM stack. This Docker image contains
an OTel collector, Prometheus (metrics), Loki (logging), Tempo (tracing) and Grafana (visualization).

```shell
docker run -p 3000:3000 -p 4317:4317 -p 4318:4318 --rm -ti grafana/otel-lgtm
```

### Running the application

We run the application like we normally do, but we set an environment variable to
use an insecure connection to the OpenTelemetry collector.

```shell
OTEL_EXPORTER_OTLP_INSECURE="true" go run cmd/app/main.go
```

In another shell we can now 'bombard' our app with some traffic. Every 5 seconds
we run our `/school/tick` endpoint, which will increment the 'schools hour'
and have ~students~ goats attend their class.

```shell
watch -n 5 curl localhost:9000/school/tick
```

## Instrumentation

We set up the OpenTelemetry Go SDK and exporters. They have defaults that are
unchanged in our LGTM container and as such we are exporting it to `localhost:4317` (HTTP)
or `localhost:4318` (gRPC). We set the `MetricProvider` and `TracerProvider` in the global space
and use the `otelhttp` package to instrument our HTTP requests.

### HTTP Middleware

There is some middleware included and wrapped around our HTTP mux. The double definition of routes
is not ideal, but in a production environment you would probably have a routing framework that
offers OTel functionality natively.

```go
// wrap our route with the route tag, to have 'http.route' in our metrics
r.Handle("GET /school/tick", otelhttp.WithRouteTag("/school/tick", http.HandlerFunc(h.schoolTick)))

// wrap the mux with OpenTelemetry
h.mux = otelhttp.NewHandler(r, "myapp")
```

The middleware wrapper is also not super fit for purpose, as it wraps our entire codebase
with the same 'operation' a.k.a. name.

However, for demonstration purposes it is sufficient.

### Tracing Spans

A big part of OpenTelemetry is the concept of spans. A span is a single operation and multiple
spans can be grouped together in a trace. A span can sit sequential in a trace, but you can also use 
the context of a 'parent' span to indicate that it's a 'child' of the parent span.

This is called [Context Propagation](https://opentelemetry.io/docs/concepts/context-propagation/) and
is a separate concept within OTel.

Because of our middleware we already have a span for our HTTP request. We can use this span to
continue our trace.

```go
func (h *Handler) schoolTick(w http.ResponseWriter, r *http.Request) {
    schoolCtx, span := h.Tracer.Start(r.Context(), "schoolTick")
    defer span.End()
}
```

In our code we have created a tracer `h.Tracer = otel.Tracer("school")` and we use this to create spans.
The context for this span is the context of the HTTP request, which has been enriched by the middleware.

A span has an 'operation' (a.k.a span name) and needs to be ended. If you have a well-defined operation then
you can just call `span.End()` at the end. If you however have multiple things going on and potentially return
early, then often we just use `defer span.End()` to ensure the span is ended at the end of the function.

In this specific case we also save the context that's being generated into `schoolCtx`. This way we can
create child spans that are linked to this span. Otherwise you can also drop the context variable with `_`.

### Metrics

Metrics are a bit more straight forward and can be made on the spot.

For example, in our `/health` route we have the following.

```go
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
```

We create a `Float64` counter and increment it by 1. If you need to update the metric in multiple places,
then you'd have to move the metric creation up to a higher level. It's also possible to have the metrics
created on a package level by using `var` and `const`.


## Dashboards

Our LTGM container comes with a Grafana instance. Visit http://localhost:3000/explore/metrics 
and click "Let's Start". This opens Grafana's metric exploration view, and you can see all metrics
with a little graph. These are simple line graphs and probably don't show the full picture. You can click 
on a metric to open it up into an overview of that metric with all its labels. This is quite handy to use
during exploration and to find all the labels and what impact they might have.

### Create the dashboard

Go to "Dashboards", click "New" and pick "Import". Paste the contents of `dashboard.json` and save.

You should now have a dashboard with some OpenTelemetry graphs (because it's a specification this
should work out of the box). It also comes with a graph for our `myapp_frontpage_total` metric and
a section for traces. Although you can also use the Explore page for traces.

(There's a data link on the traces table, if you click a trace ID you can click "Follow Trace ID"
to have it show up in the graph above)

You are allowed, and should, edit this dashboard and add your own metrics.

## Assignment

Alright, it's time to play with this code. Maybe you can..

- add a new service and have the traces go a level deeper
- add a service for class grades and create metrics for students grades
- simulate class engagement (0-100) by using a Gauge meter
- implement a heatmap graph in the dashboard specifically for the `/school/tick` route

Some more information on instrumenting within Go can be found at https://opentelemetry.io/docs/languages/go/instrumentation/
