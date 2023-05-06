package middleware

import (
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
	"github.com/viettranx/service-context/core"
)

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	//"Authorization" : "Bearer {token}"

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", core.ErrUnauthorized.WithError("missing access token")
	}

	return parts[1], nil
}

func RequireAuth(client clerk.Client) func(*gin.Context) {
	return func(c *gin.Context) {
		token, _ := extractTokenFromHeaderString(c.GetHeader("Authorization"))

		sessClaims, err := client.VerifyToken(token)
		if err != nil {
			panic(core.ErrUnauthorized.WithError(err.Error()))
		}

		user, err := client.Users().Read(sessClaims.Claims.Subject)
		if err != nil {
			panic(core.ErrUnauthorized.WithError(err.Error()))
		}

		c.Set(core.KeyRequester, user)
		c.Next()
	}
}
