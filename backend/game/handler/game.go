package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"skyhawk/backend/game/domain"
	"skyhawk/backend/game/usecase"
)

type Handler struct {
	useCase *usecase.UseCase
	logger  *zap.Logger
}

func NewHandler(useCase *usecase.UseCase, logger *zap.Logger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

func (h *Handler) GameLogHandler(c echo.Context) error {
	var req domain.GameStatsReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := h.useCase.LogGame(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

func (h *Handler) TeamSeasonStatsHandler(c echo.Context) error {
	id := c.Param("id")

	stats, err := h.useCase.GetTeamSeasonStats(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) GameStatsHandler(c echo.Context) error {
	id := c.Param("id")

	res, err := h.useCase.GetGameStats(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) PlayerSeasonStatsHandler(c echo.Context) error {
	playerId := c.Param("player_id")

	result, err := h.useCase.GetPlayerSeasonStats(playerId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	return c.JSON(http.StatusOK, result)
}
