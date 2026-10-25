// Package metrics owns the Prometheus registry of the service:
// Go runtime / process collectors, HTTP metrics and business metrics.
// They are exposed in OpenMetrics format on a dedicated port (METRICS_ADDR),
// separate from the public API port.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

const namespace = "meeting"

// Label values used by business metrics. Kept as constants so that every
// series is pre-created with 0 and is visible in VictoriaMetrics right away.
const (
	ResultSuccess = "success"
	ResultFailure = "failure"

	LoginMethodPassword = "password"
	LoginMethodDummy    = "dummy"

	BookingFailSlotAlreadyBooked = "slot_already_booked"
	BookingFailSlotInPast        = "slot_in_past"
	BookingFailSlotNotFound      = "slot_not_found"
	BookingFailForbidden         = "forbidden"
	BookingFailInvalidRequest    = "invalid_request"
	BookingFailInternal          = "internal"

	ConfLinkSuccess     = "success"
	ConfLinkUnavailable = "unavailable"
	ConfLinkInternal    = "internal_error"
)

type Metrics struct {
	Registry *prometheus.Registry

	// --- framework / HTTP ---
	HTTPRequests        *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec

	// --- business ---
	BookingsCreated   *prometheus.CounterVec
	BookingsCancelled prometheus.Counter
	BookingFailures   *prometheus.CounterVec
	Logins            *prometheus.CounterVec
	UsersRegistered   *prometheus.CounterVec
	RoomsCreated      prometheus.Counter
	ConfLinkRequests  *prometheus.CounterVec
}

func New() *Metrics {
	reg := prometheus.NewRegistry()

	m := &Metrics{
		Registry: reg,

		HTTPRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_requests_total",
			Help:      "Number of HTTP requests by method, route pattern and status code.",
		}, []string{"method", "route", "status"}),

		HTTPRequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latency by method and route pattern.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"method", "route"}),

		BookingsCreated: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "bookings_created_total",
			Help:      "Number of successfully created bookings.",
		}, []string{"with_conference_link"}),

		BookingsCancelled: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "bookings_cancelled_total",
			Help:      "Number of active bookings cancelled by users.",
		}),

		BookingFailures: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "booking_failures_total",
			Help:      "Number of rejected booking attempts by reason.",
		}, []string{"reason"}),

		Logins: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "auth_logins_total",
			Help:      "Number of login attempts by method and result.",
		}, []string{"method", "result"}),

		UsersRegistered: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "users_registered_total",
			Help:      "Number of successfully registered users by role.",
		}, []string{"role"}),

		RoomsCreated: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "rooms_created_total",
			Help:      "Number of meeting rooms created by admins.",
		}),

		ConfLinkRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "conference_link_requests_total",
			Help:      "Calls to the external conference service by result.",
		}, []string{"result"}),
	}

	reg.MustRegister(
		// runtime: goroutines, GC, heap, memstats, scheduler (go_*)
		collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsGC, collectors.MetricsMemory, collectors.MetricsScheduler),
		),
		// process: CPU, RSS, open fds, start time (process_*)
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),

		m.HTTPRequests,
		m.HTTPRequestDuration,
		m.BookingsCreated,
		m.BookingsCancelled,
		m.BookingFailures,
		m.Logins,
		m.UsersRegistered,
		m.RoomsCreated,
		m.ConfLinkRequests,
	)

	m.initSeries()

	return m
}

// initSeries creates every known label combination with value 0,
// so rate()/increase() work from the very first scrape.
func (m *Metrics) initSeries() {
	m.BookingsCreated.WithLabelValues("true")
	m.BookingsCreated.WithLabelValues("false")

	for _, r := range []string{
		BookingFailSlotAlreadyBooked, BookingFailSlotInPast, BookingFailSlotNotFound,
		BookingFailForbidden, BookingFailInvalidRequest, BookingFailInternal,
	} {
		m.BookingFailures.WithLabelValues(r)
	}

	for _, method := range []string{LoginMethodPassword, LoginMethodDummy} {
		m.Logins.WithLabelValues(method, ResultSuccess)
		m.Logins.WithLabelValues(method, ResultFailure)
	}

	m.UsersRegistered.WithLabelValues("user")
	m.UsersRegistered.WithLabelValues("admin")

	for _, r := range []string{ConfLinkSuccess, ConfLinkUnavailable, ConfLinkInternal} {
		m.ConfLinkRequests.WithLabelValues(r)
	}
}
