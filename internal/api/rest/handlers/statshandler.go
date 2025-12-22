package handlers

import (
	"biturl/internal/api/rest"
	ratelimiter "biturl/internal/middleware/rate-limiter"
	"biturl/internal/repository"
	"biturl/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type StatsHander struct {
	StatsSVC  service.StatsSVC
	RateLimit ratelimiter.RateLimiter
}

func SetupStatsRoute(rh *rest.RestHandler) {
	app := rh.App

	statsSVC := service.StatsSVC{
		ClickhouseRepo: *repository.NewClkHouseRepo(rh.ClickhouseConn),
		PGRepo:         repository.NewPostgresRepo(rh.DB),
	}

	handler := StatsHander{
		StatsSVC:  statsSVC,
		RateLimit: *rh.StatsRatelimit,
	}

	statsRoutes := app.Group("/api/stats", rh.StatsRatelimit.Middleware())

	statsRoutes.Get("/:stats_access_key", handler.GetStats)
}

func (h *StatsHander) GetStats(ctx *fiber.Ctx) error {

	statsAccessKey := ctx.Params("stats_access_key")

	c := ctx.UserContext()

	stats, err := h.StatsSVC.GetStats(c, statsAccessKey)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get stats",
			"detail":  err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "stats loaded successfully",
		"data":    stats,
	})
}
