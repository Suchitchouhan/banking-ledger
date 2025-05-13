// internal/api/router.go
package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Account routes
	r.POST("/accounts", h.CreateAccountHandler)
	r.GET("/accounts", h.ListAccountsHandler)
	r.GET("/accounts/:id", h.GetAccountHandler)

	// Transaction routes
	r.POST("/accounts/:id/deposit", h.DepositHandler)
	r.POST("/accounts/:id/withdraw", h.WithdrawHandler)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	return r
}

// StartServer starts the HTTP server
func StartServer(router *gin.Engine, port string) error {
	log.Printf("Starting server on port %s", port)
	return router.Run(":" + port)
}
