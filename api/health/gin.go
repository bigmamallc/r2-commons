package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GinHandler(ctx *gin.Context) {
	status := GatherStatus()
	code := http.StatusOK
	if !status.Healthy {
		code = http.StatusInternalServerError
	}

	ctx.JSON(code, status)
}
