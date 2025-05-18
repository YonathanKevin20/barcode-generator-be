package middleware

import (
	"barcode-generator-be/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token format")
	}

	claims, err := utils.GetTokenClaims(tokenString)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token")
	}

	isBlacklisted, err := utils.TokenBlacklist.IsBlacklisted(jti)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error checking token status")
	}
	if isBlacklisted {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token has been invalidated")
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid ID in token")
	}
	id := uint(idFloat)

	username, ok := claims["username"].(string)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid username in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid role in token")
	}

	c.Locals("id", id)
	c.Locals("username", username)
	c.Locals("role", role)
	c.Locals("jti", jti)

	return c.Next()
}
