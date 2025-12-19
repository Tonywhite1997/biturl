package handlers

import (
	"biturl/internal/api/rest"
	"biturl/internal/dto"
	"biturl/internal/helper"
	"biturl/internal/helper/geo"
	ratelimiter "biturl/internal/middleware/rate-limiter"
	"biturl/internal/repository"
	"biturl/internal/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mssola/user_agent"
)

type URLhandler struct {
	Svc          service.URLsvc
	GEO          *geo.GeoRedisCache
	URLRateLimit *ratelimiter.RateLimiter
}

func SetupURLroutes(rh *rest.RestHandler) {

	app := rh.App
	svc := service.URLsvc{
		PG:           repository.NewPostgresRepo(rh.DB),
		RDB:          repository.NewRedisRepo(rh.RDB),
		ClkhouseConn: *repository.NewClkHouseRepo(rh.ClickhouseConn),
		RabbitConn:   rh.RabbitConn,
	}

	geo := rh.GEODB

	handler := URLhandler{
		Svc:          svc,
		GEO:          geo,
		URLRateLimit: rh.URLRateLimit,
	}

	app.Get("/:shortcode", handler.URLRateLimit.Middleware(), handler.LoadURL)

	urlRoutes := app.Group("/api/url", handler.URLRateLimit.Middleware())
	urlRoutes.Post("/shorten", handler.CreateShortURL)
	urlRoutes.Delete("/:shortcode", handler.DeleteURL)
	urlRoutes.Patch("/:accesskey", handler.IncreaseExpiryDate)

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

	shortCode, stats_access_key, err := h.Svc.CreateShortURL(req, c)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":  "error occured",
			"detail": err.Error(),
		})
	}

	baseurl := ctx.BaseURL()

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"shorturl":       baseurl + "/" + shortCode,
		"statsAccessKey": stats_access_key,
	})

}

func (h *URLhandler) LoadURL(ctx *fiber.Ctx) error {
	shortCode := ctx.Params("shortcode")

	ua := ctx.Get("User-Agent")
	uaParser := user_agent.New(ua)
	browser, _ := uaParser.Browser()
	device := uaParser.Platform()
	os := uaParser.OS()
	ip := ctx.IP()
	// ip := "8.8.8.8"

	c := ctx.UserContext()
	country, city, err := h.GEO.LookupIP(ip, c)
	if err != nil {
		fmt.Printf("cannot get country: %v", err)
	}

	stats := repository.Stats{
		Id:           helper.GenerateShortCode(),
		Url_short_id: shortCode,
		User_ip:      ip,
		User_agent:   ua,
		Referer:      ctx.Get("referer"),
		Device:       device,
		OS:           os,
		Browser:      browser,
		Country:      country,
		City:         city,
		Timestamp:    time.Now(),
	}

	c = ctx.UserContext()
	url, err := h.Svc.LoadURL(shortCode, c, stats)

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

func (h *URLhandler) IncreaseExpiryDate(ctx *fiber.Ctx) error {
	accessKey := ctx.Params("accesskey")

	err := h.Svc.IncreaseExpiryDate(accessKey)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "failed to increase expiry date",
			"details": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "expiry date updated",
	})

}
