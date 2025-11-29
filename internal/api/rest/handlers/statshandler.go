package handlers

import (
	"biturl/internal/api/rest"
	"biturl/internal/repository"
	"biturl/internal/service"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type StatsHander struct {
	StatsSVC service.StatsSVC
}

func SetupStatsRoute(rh *rest.RestHandler) {
	app := rh.App

	statsSVC := service.StatsSVC{
		ClickhouseRepo: *repository.NewClkHouseRepo(rh.ClickhouseConn),
		PGRepo:         repository.NewPostgresRepo(rh.DB),
	}

	handler := StatsHander{
		StatsSVC: statsSVC,
	}

	statsRoutes := app.Group("/stats")

	statsRoutes.Get("/:stats_access_key", handler.GetStatsByShortCode)
}

func (h *StatsHander) GetStatsByShortCode(ctx *fiber.Ctx) error {

	statsAccessKey := ctx.Params("stats_access_key")

	fmt.Println(statsAccessKey)

	c := ctx.UserContext()

	stats, err := h.StatsSVC.GetStatsByShortCode(c, statsAccessKey)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get stats",
			"details": err.Error(),
		})
	}

	if len(stats) == 0 {
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "no stats found for this URL",
			"data":    []repository.Stats{},
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "stats loaded successfully",
		"data":    stats,
	})
}
