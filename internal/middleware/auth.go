package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "payment-gateway/pkg/response"
)

func AuthMiddleware(apiKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.GetHeader("X-API-Key")
        
        if key == "" {
            response.Error(c, http.StatusUnauthorized, "missing API key")
            c.Abort()
            return
        }

        if key != apiKey {
            response.Error(c, http.StatusUnauthorized, "invalid API key")
            c.Abort()
            return
        }

        c.Next()
    }
}