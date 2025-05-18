package utils

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	Data      any   `json:"data"`
	Total     int64 `json:"total"`
	Page      int   `json:"page"`
	TotalPage int   `json:"total_page"`
	Limit     int   `json:"limit"`
}

func JSONResponse(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"data": data,
	})
}

func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": fiber.Map{
			"message": message,
		},
	})
}

func PaginatedResponse(c *fiber.Ctx, data any, total int64, page int, limit int) error {
	response := Pagination{
		Data:      data,
		Total:     total,
		Page:      page,
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Limit:     limit,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
