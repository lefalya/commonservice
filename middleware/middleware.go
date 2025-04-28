package middleware

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lefalya/commonservice/jwt"
	"log/slog"
	"math/rand"
	"os"
	"strings"
)

var (
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

type ErrorResponse struct {
	Code string `json:"code"`
	ID   string `json:"id"`
}

func ConstructErrorResponse(c *fiber.Ctx, component string, status int, error error, code string, inputBody string, source string) error {
	errorId := RandId(10)

	response := ErrorResponse{
		Code: code,
		ID:   errorId,
	}

	Logger.Error("endpoint-error", "component", component, "source", source, "code", code, "error", error.Error(), "ID", errorId, "input", inputBody)

	c.Set("Content-Type", "application/json")
	return c.Status(status).JSON(response)
}

func RandId(length int) string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}

	return string(result)
}

func JWTMiddleware(componentName string, mandatory bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			if !mandatory {
				return c.Next()
			}
			return ConstructErrorResponse(c, componentName, fiber.StatusUnauthorized, errors.New("Missing JWT token in the header"), "MX401", "", "JWTMiddleware")
		}

		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) < 2 {
			return ConstructErrorResponse(c, componentName, fiber.StatusUnauthorized, errors.New("Missing JWT token in the header"), "MX401", "", "JWTMiddleware")
		}

		claims, err := jwt.JWTDecode(tokenString[1])
		if err != nil {
			return ConstructErrorResponse(c, componentName, fiber.StatusUnauthorized, errors.New("Invalid JWT token"), "MX401", "", "JWTMiddleware")
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}

func ParseCredential(c *fiber.Ctx) *jwt.Claims {
	var claim *jwt.Claims

	rawClaims := c.Locals("claims")
	if rawClaims == nil {
		return claim
	}

	claim, ok := rawClaims.(*jwt.Claims)
	if !ok {
		return claim
	}

	return claim
}

func StringifyBody[T any](s T) string {
	data, err := json.Marshal(s)
	if err != nil {
		return "{}"
	}

	return string(data)
}
