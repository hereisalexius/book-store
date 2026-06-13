package handler

import (
	"context"
	"database/sql"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db     *sql.DB
	commit string
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	commit := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				commit = s.Value[:7]
				break
			}
		}
	}
	return &HealthHandler{db: db, commit: commit}
}

// Health godoc
// @Summary      Health check
// @Description  Returns service status, commit hash, and DB connectivity.
// @Tags         actuator
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      503  {object}  map[string]string
// @Router       /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "ok"
	status := http.StatusOK
	if err := h.db.PingContext(ctx); err != nil {
		dbStatus = "error"
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"status": map[int]string{http.StatusOK: "ok", http.StatusServiceUnavailable: "degraded"}[status],
		"commit": h.commit,
		"db":     dbStatus,
	})
}
