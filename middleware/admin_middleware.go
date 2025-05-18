package middleware

import (
	"barcode-generator-be/models"
	"barcode-generator-be/utils"

	"github.com/gofiber/fiber/v2"
)

func AdminOnly(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != string(models.RoleAdmin) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Admin access required")
	}

	return c.Next()
}
