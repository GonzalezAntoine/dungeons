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

// Get controller to get list of dungeon
func (s *Dungeon) Get(ctx *gin.Context) {
	var params models.QueryParams

	params.Parse(ctx)
	messageTypes := &models.MessageTypes{
		OK:                  "dungeon.Search.Found",
		BadRequest:          "dungeon.Search.BadRequest",
		NotFound:            "dungeon.Search.NotFound",
		InternalServerError: "dungeon.Search.Error",
	}

	dungeons, err := s.DungeonService.Get(params)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	totalCount := len(dungeons)
	if totalCount == 0 {
		status := http.StatusNotFound
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, err))
		return
	}

	low := params.Offset - 1
	if low == -1 {
		low = 0
	}

	// Available CountMax calculation
	maxCount := params.Count
	if maxCount == 0 {
		maxCount = 100
	}

	high := maxCount + low
	if high > totalCount {
		high = totalCount
	}

	if low > high {
		status := http.StatusBadRequest
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.BadRequest, err))
		return
	}

	sendingDungeons := dungeons[low:high]
	meta := models.MetaResponse{
		ObjectName: "Dungeon",
		TotalCount: totalCount,
		Count:      len(sendingDungeons),
		Offset:     low + 1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: sendingDungeons,
	}
	common.SendResponse(ctx, http.StatusOK, response)
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

func (s *Dungeon) GetByID(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "dungeon.GetByID.Found",
		NotFound:            "dungeon.GetByID.NotFound",
		BadRequest:          "dungeon.GetByID.BadRequest",
		InternalServerError: "dungeon.GetByID.Error",
	}

	id := ctx.Param("id")

	dungeon, err := s.DungeonService.GetByID(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	// Créer la réponse avec le client
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

	// Envoyer la réponse
	common.SendResponse(ctx, http.StatusOK, response)
}

func (s *Dungeon) Update(ctx *gin.Context) {
	var in models.Dungeon
	messageTypes := &models.MessageTypes{
		OK:                  "dungeon.Update.Updated",
		BadRequest:          "dungeon.Update.BadRequest",
		InternalServerError: "dungeon.Update.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	id := ctx.Param("id")
	err := s.DungeonService.Update(id, &in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, messageTypes.OK, "dungeon updated"))
}

func (s *Dungeon) Status(ctx *gin.Context) {
	id := ctx.Param("id")
	err := s.DungeonService.Publish(id, nil)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, "dungeon.Publish.Error", err))
		return
	}
	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "dungeon.Status.Updated", "dungeon status updated"))
}