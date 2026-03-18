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
		dungeons := v1.Group("/dungeons")
		{
			dungeons.GET("", dungeonController.Get)
			dungeons.POST("", dungeonController.Create)
		}
	}
}