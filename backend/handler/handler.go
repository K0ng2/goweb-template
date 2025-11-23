package handler

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"

	"project/database"
	"project/model"
	"project/repo"
)

type Handler struct {
	repo repo.Repository
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{
		repo: repo.NewRepo(db),
	}
}

var startTime = time.Now()

// DatabaseHealth godoc
// @Summary Database Health check
// @Description Check the Database health status
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} model.DatabaseHealth
// @Router /databasez [get]
func (h *Handler) DatabaseHealth(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	ping := h.repo.Ping(ctx)
	// Check database connectivity
	dbStatus := "healthy"
	if err := ping; err != nil {
		dbStatus = "unhealthy"
	}

	// Calculate uptime
	uptime := time.Since(startTime).String()

	response := model.DatabaseHealth{
		Status:    "healthy",
		Timestamp: time.Now(),
		Database:  dbStatus,
		Uptime:    uptime,
	}

	// If database is unhealthy, set overall status to unhealthy
	if dbStatus == "unhealthy" {
		response.Status = "degraded"
		return c.Status(http.StatusServiceUnavailable).JSON(response)
	}

	log.Info("test")
	return c.JSON(response)
}

func Response[T any](r T, m *model.Meta) model.APIResponse[T] {
	return model.APIResponse[T]{
		Data: r,
		Meta: m,
	}
}

func Error(err error) model.APIResponse[any] {
	return model.APIResponse[any]{
		Error: err.Error(),
	}
}

// GetOffset extracts limit and offset query parameters from the request context.
func GetOffset(c fiber.Ctx) (*model.Offset, error) {
	var query model.Offset

	if err := c.Bind().Query(&query); err != nil {
		return nil, err
	}

	return &query, nil
}
