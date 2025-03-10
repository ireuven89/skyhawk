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

func (h *Handler) GameStatsHandler(c echo.Context) error {
	var req domain.GameStats

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

	stats, err := h.useCase.FindTeam(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) TeamStatsStatsHandler(c echo.Context) error {
	var req domain.GameStats

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := h.useCase.LogGame(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

func (h *Handler) PlayerStatsHandler(c echo.Context) error {
	var req domain.GameStats

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := h.useCase.LogGame(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}
func (h *Handler) PlayerSeasonStatsHandler(c echo.Context) error {
	var req domain.GameStats

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := h.useCase.LogGame(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}
