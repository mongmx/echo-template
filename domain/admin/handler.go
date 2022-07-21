package admin

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	store *Store
	rd    *redis.Client
}

func NewHandler(store *Store, rd *redis.Client) *Handler {
	return &Handler{store: store, rd: rd}
}

func (h *Handler) ListProfiles(c echo.Context) error {
	us, err := h.store.ListUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, newListProfileResponse(us))
}
