package middlewares

import (
	"docmate/config"
	"docmate/internal/model"
	"docmate/response"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// AuthRoles middleware verifies the JWT token and checks if the user has one of the allowed roles.
func AuthRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.JSON(c, http.StatusUnauthorized, false, "Missing Authorization header", nil, nil)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.JSON(c, http.StatusUnauthorized, false, "Invalid Authorization header format", nil, nil)
			}
			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(config.App().JWTSecret), nil
			})

			if err != nil || !token.Valid {
				return response.JSON(c, http.StatusUnauthorized, false, "Invalid or expired token", nil, nil)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return response.JSON(c, http.StatusUnauthorized, false, "Invalid token claims", nil, nil)
			}

			role, ok := claims["role"].(string)
			if !ok {
				return response.JSON(c, http.StatusForbidden, false, "Access denied", nil, nil)
			}

			// Check if the user's role is in the allowedRoles list
			roleAllowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					roleAllowed = true

					break
				}
			}

			if !roleAllowed {
				return response.JSON(c, http.StatusForbidden, false, "Access denied", nil, nil)
			}

			headers := c.Request().Header
			id, _ := strconv.Atoi(headers.Get(headerUserID))

			user := &model.User{
				ID:       id,
				UserName: headers.Get(headerUserName),
				Email:    headers.Get(headerUserEmail),
				Role:     role,
			}

			// Store user details in context
			c.Set("user", user)

			return next(c)
		}
	}
}
