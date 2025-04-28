package middleware

import (
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

func RandId(length int) string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}

	return string(result)
}

type ErrorResponse struct {
	Code string `json:"code"`
	ID   string `json:"id"`
}

type Middleware struct {
	componentName string
}

func (m Middleware) ConstructErrorResponse(c *fiber.Ctx, status int, error error, code string, source string) error {
	errorId := RandId(10)

	response := ErrorResponse{
		Code: code,
		ID:   errorId,
	}

	var inputBody string
	if c.Request().Body() != nil && len(c.Request().Body()) > 0 {
		inputBody = string(c.Request().Body())
	}

	Logger.Error("endpoint-error", "component", m.componentName, "source", source, "code", code, "error", error.Error(), "ID", errorId, "input", inputBody)

	c.Set("Content-Type", "application/json")
	return c.Status(status).JSON(response)
}

func (m Middleware) ValidateJWT(mandatory bool, jwtDecode func(string) (*jwt.Claims, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			if !mandatory {
				return c.Next()
			}
			return m.ConstructErrorResponse(c, fiber.StatusUnauthorized, errors.New("Missing JWT token in the header"), "MX401", "JWTMiddleware")
		}

		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) < 2 {
			return m.ConstructErrorResponse(c, fiber.StatusUnauthorized, errors.New("Missing JWT token in the header"), "MX401", "JWTMiddleware")
		}

		claims, err := jwtDecode(tokenString[1])
		if err != nil {
			return m.ConstructErrorResponse(c, fiber.StatusUnauthorized, errors.New("Invalid JWT token"), "MX401", "JWTMiddleware")
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}

func (m Middleware) ParseCredential(c *fiber.Ctx) *jwt.Claims {
	var claim *jwt.Claims

	rawClaims := c.Locals("claims")
	if rawClaims == nil {
		return claim
	}

	claim, _ = rawClaims.(*jwt.Claims)
	return claim
}

func New(componentName string) Middleware {
	return Middleware{
		componentName: componentName,
	}
}
