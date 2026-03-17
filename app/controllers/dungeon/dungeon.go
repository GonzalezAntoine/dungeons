package dungeon

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	dungeon "dungeons/app/services/dungeons"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Dungeon struct {
	DungeonService *dungeon.Dungeon
}

func New(dungeonService *dungeon.Dungeon) *Dungeon {
	return &Dungeon{
		DungeonService: dungeonService,
	}
}

func (s *Dungeon) Create(ctx *gin.Context) {
	var in models.Dungeon
	
	messageTypes := &models.MessageTypes{
		Created:                  "dungeon.Create.Created",
		BadRequest:          "dungeon.Create.BadRequest",
		InternalServerError: "dungeon.Create.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	dungeon, err := s.DungeonService.Create(&in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	meta := models.MetaResponse{
		ObjectName: "Dungeon",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: dungeon,
	}
	common.SendResponse(ctx, http.StatusCreated, response)
}