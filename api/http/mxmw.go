package http

import (
	"github.com/gin-gonic/gin"
	. "github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type metrics struct {
	apiResponses    *CounterVec
	apiResponseTime *HistogramVec
}

const labelStatus = "status"
const labelURI = "uri"

func RequestMetricsMW(r *Registry, subsystem string, responseTimeBucketsSec []float64) gin.HandlerFunc {
	mx := &metrics{
		apiResponses: NewCounterVec(CounterOpts{Name: "apiResponses", Subsystem: subsystem},
			[]string{labelURI, labelStatus}),
		apiResponseTime: NewHistogramVec(HistogramOpts{Name: "apiResponseTime",
			Subsystem: subsystem, Buckets: responseTimeBucketsSec}, []string{labelURI, labelStatus}),
	}
	r.MustRegister(mx.apiResponses, mx.apiResponseTime)

	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		dur := time.Now().Sub(start)

		uri := ctx.FullPath()
		status := strconv.Itoa(ctx.Writer.Status())
		labels := Labels{labelURI: uri, labelStatus: status}

		mx.apiResponses.With(labels).Inc()
		mx.apiResponseTime.With(labels).Observe(dur.Seconds())
	}
}
