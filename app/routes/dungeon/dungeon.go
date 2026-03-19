package dungeon

import (
	controller "dungeons/app/controllers/dungeon"
	service "dungeons/app/services/dungeons"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {

	serviceDungeon := service.New()
	dungeonController := controller.New(serviceDungeon)

	v1 := g.Group("/v1")
	{
		mj := v1.Group("/mj")
		{
			dungeons := mj.Group("/dungeons")
			{
				dungeons.POST("", dungeonController.Create)
				dungeons.POST("/:id/publish", dungeonController.Status)
				//dungeons.POST("/:id/steps", dungeonController.UpdateSteps)
			}
		}
	}
}