package handlers

import (
	"barcode-generator-be/models"
	"barcode-generator-be/utils"

	"github.com/gofiber/fiber/v2"
)

var statusRepo models.StatusRepository = &models.GormStatusRepository{}

func SetStatusRepository(repo models.StatusRepository) {
	statusRepo = repo
}

func GetStatuses(c *fiber.Ctx) error {
	statuses, err := statusRepo.FindAll()
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, fiber.Map{"error": "Failed to fetch statuses"})
	}
	return utils.JSONResponse(c, fiber.StatusOK, statuses)
}
