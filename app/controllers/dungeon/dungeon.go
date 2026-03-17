package dungeon

import (
	"dungeons/app/controllers/common"
	"dungeons/app/controllers/dungeon"
	"dungeons/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Dungeon struct {
	DungeonService *dungeon.Dungeon
}

type BossStep struct {
	//
}

func New(dungeonService *dungeon.Dungeon) *Dungeon {
	return &Dungeon{
		DungeonService: dungeonService,
	}
}

func (s *Dungeon) AddBossStep(bossStep BossStep) {
	// Implementation for adding boss step
}

func (s *Dungeon) Create(ctx *gin.Context) {
	var in models.Dungeon
	
	messageTypes := &models.MessageTypes{
		OK:                  "dungeon.Create.Created",
		BadRequest:          "dungeon.Create.BadRequest",
		NotFound:            "dungeon.Create.NotFound",
		InternalServerError: "dungeon.Create.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	dungeon, err := s.DungeonService.Create(in)
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