package apiroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/kanhaiyagupta9045/kirana_club/internals/service"
)

func StoreVisitServiceRoutes(router *gin.Engine) {
	router.POST("/api/submit", service.SubmitJobHandler())
	router.GET("/api/status", service.GetJobInfoHandler())
}
