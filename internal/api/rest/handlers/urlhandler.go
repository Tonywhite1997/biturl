package handlers

import (
	"biturl/internal/api/rest"
	"biturl/internal/dto"
	"biturl/internal/repository"
	"biturl/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type URLhandler struct {
	Svc service.URLsvc
}

func SetupURLroutes(rh *rest.RestHandler) {

	app := rh.App
	svc := service.URLsvc{
		PG:         repository.NewPostgresRepo(rh.DB),
		RDB:        repository.NewRedisRepo(rh.RDB),
		RabbitConn: rh.RabbitConn,
	}

	handler := URLhandler{
		Svc: svc,
	}

	app.Post("/", handler.CreateShortURL)
	app.Get("/:shortcode", handler.LoadURL)
	app.Delete("/:shortcode", handler.DeleteURL)

}

func (h *URLhandler) CreateShortURL(ctx *fiber.Ctx) error {
	req := dto.URLdto{}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	c := ctx.UserContext()

	shortCode, err := h.Svc.CreateShortURL(req, c)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "error occured",
			"detail": err.Error(),
		})
	}

	baseurl := ctx.BaseURL()

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"shorturl": baseurl + "/" + shortCode,
		"status":   "ok",
	})

}

func (h *URLhandler) LoadURL(ctx *fiber.Ctx) error {
	shortCode := ctx.Params("shortcode")

	c := ctx.UserContext()
	url, err := h.Svc.LoadURL(shortCode, c)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "an error occured",
			"detail": err.Error(),
		})
	}
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

func (h *URLhandler) DeleteURL(ctx *fiber.Ctx) error {
	shortCode := ctx.Params("shortcode")

	c := ctx.UserContext()
	err := h.Svc.DeleteURL(shortCode, c)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "error occured",
			"detail": err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "url deleted",
	})
}
