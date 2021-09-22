package http

import (
	"github.com/gin-gonic/gin"
	. "github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type metrics struct {
	responses    *CounterVec
	responseTime *HistogramVec
}

const labelStatus = "status"
const labelURI = "uri"

func RequestMetricsMW(r *Registry, subsystem string, responseTimeBucketsSec []float64) gin.HandlerFunc {
	mx := &metrics{
		responses: NewCounterVec(CounterOpts{Name: "responses", Subsystem: subsystem},
			[]string{labelURI, labelStatus}),
		responseTime: NewHistogramVec(HistogramOpts{Name: "responseTime",
			Subsystem: subsystem, Buckets: responseTimeBucketsSec}, []string{labelURI, labelStatus}),
	}
	r.MustRegister(mx.responses, mx.responseTime)

	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		dur := time.Now().Sub(start)

		uri := ctx.FullPath()
		status := strconv.Itoa(ctx.Writer.Status())
		labels := Labels{labelURI: uri, labelStatus: status}

		mx.responses.With(labels).Inc()
		mx.responseTime.With(labels).Observe(dur.Seconds())
	}
}
