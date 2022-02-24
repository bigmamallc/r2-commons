package http

import (
	"github.com/gin-gonic/gin"
	"net/http/pprof"
)


// Adds the pprof endpoints
// Note that they will be unauthenticated by default!
func AddPprofEndpoints(eng *gin.Engine) *gin.RouterGroup {
	pprofHandler := gin.WrapF(pprof.Index)
	debugGroup := eng.Group("/debug")
	debugGroup.Any("/pprof/", pprofHandler)
	debugGroup.Any("/pprof/allocs", pprofHandler)
	debugGroup.Any("/pprof/block", pprofHandler)
	debugGroup.Any("/pprof/goroutine", pprofHandler)
	debugGroup.Any("/pprof/heap", pprofHandler)
	debugGroup.Any("/pprof/mutex", pprofHandler)
	debugGroup.Any("/pprof/threadcreate", pprofHandler)
	debugGroup.Any("/pprof/cmdline", gin.WrapF(pprof.Cmdline))
	debugGroup.Any("/pprof/profile", gin.WrapF(pprof.Profile))
	debugGroup.Any("/pprof/symbol", gin.WrapF(pprof.Symbol))
	debugGroup.Any("/pprof/trace", gin.WrapF(pprof.Trace))
	return debugGroup
}
