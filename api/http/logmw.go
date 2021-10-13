package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"regexp"
	"time"
)

func RequestLogMW(log zerolog.Logger, skipRegex *regexp.Regexp) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if skipRegex != nil && skipRegex.MatchString(ctx.Request.RequestURI) {
			ctx.Next()
			return
		}

		start := time.Now()
		ctx.Next()
		dur := time.Now().Sub(start)

		ev := log.Debug()
		if ctx.Writer.Status() >= 500 {
			ev = log.Warn()
		} else if ctx.Request.Method != http.MethodGet || (ctx.Writer.Status() < 200 || ctx.Writer.Status() >= 300) {
			ev = log.Info()
		}
		ev.
			Str("method", ctx.Request.Method).
			Str("uri", ctx.Request.RequestURI).
			Int("status", ctx.Writer.Status()).
			Dur("dur", dur).
			Msg("request processed")
	}
}
