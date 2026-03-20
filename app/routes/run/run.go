package run

import (
	controller "dungeons/app/controllers/run"
	service "dungeons/app/services/runs"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {

	serviceRun := service.New()
	runController := controller.New(serviceRun)

	v1 := g.Group("/v1")
	{
		run := v1.Group("/runs")
		{
			run.GET("", runController.Get)
			run.GET("/:id", runController.GetByID)
			run.POST("", runController.Create)
		}
	}
}