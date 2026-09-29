package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/naseej/server/config"
	"github.com/naseej/server/internal/graph"
	"github.com/naseej/server/internal/middleware"
	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/internal/websocket"
	"github.com/naseej/server/pkg/database"

	"github.com/google/uuid"
	gorilla "github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func main() {
	// Configure structured logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load .env file (optional in production)
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	// Connect to database
	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	log.Info().Msg("connected to database")

	// Run AutoMigrate
	if err := autoMigrate(db); err != nil {
		log.Fatal().Err(err).Msg("failed to auto-migrate database")
	}
	log.Info().Msg("database migration complete")

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()

	// Initialize resolver (dependency injection)
	resolver := graph.NewResolver(db, cfg, wsHub)

	// Setup Gin
	gin.SetMode(cfg.Server.GinMode)
	r := gin.New()

	// Global middleware
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// GraphQL endpoint (requires tenant context)
	api := r.Group("/")
	api.Use(middleware.TenantMiddleware(db))
	{
		api.POST("/graphql", graphqlHandler(resolver, cfg))
		api.GET("/graphql", playgroundHandler())
	}

	// Auth endpoints (require tenant but not auth)
	auth := r.Group("/auth")
	auth.Use(middleware.TenantMiddleware(db))
	{
		auth.POST("/login", loginHandler(resolver))
		auth.POST("/refresh", refreshHandler(resolver))
	}

	// WebSocket endpoint
	r.GET("/ws", wsHandler(wsHub, db))

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Info().Str("addr", addr).Msg("starting server")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Merchant{},
		&models.Location{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.ManagerOverride{},
		&models.SizeGroup{},
		&models.Size{},
		&models.Category{},
		&models.Color{},
		&models.Brand{},
		&models.Product{},
		&models.ProductImage{},
		&models.ProductVariant{},
		&models.InventoryLevel{},
		&models.InventoryMovement{},
		&models.Customer{},
		&models.CustomerAddress{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Refund{},
		&models.Coupon{},
		&models.LoyaltyTransaction{},
	)
}

func graphqlHandler(resolver *graph.Resolver, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, return a simple response indicating GraphQL is ready
		// Full gqlgen integration will be generated with `go run github.com/99designs/gqlgen generate`
		c.JSON(http.StatusOK, gin.H{
			"message": "GraphQL endpoint ready. Run 'go run github.com/99designs/gqlgen generate' to generate the runtime.",
		})
	}
}

func playgroundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, playgroundHTML)
	}
}

func loginHandler(resolver *graph.Resolver) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
			return
		}

		payload, err := resolver.AuthService.Login(c.Request.Context(), input.Email, input.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":         payload.Token,
			"refresh_token": payload.RefreshToken,
			"user": gin.H{
				"id":         payload.User.ID,
				"email":      payload.User.Email,
				"first_name": payload.User.FirstName,
				"last_name":  payload.User.LastName,
			},
		})
	}
}

func refreshHandler(resolver *graph.Resolver) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Token string `json:"token" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
			return
		}

		payload, err := resolver.AuthService.RefreshToken(c.Request.Context(), input.Token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":         payload.Token,
			"refresh_token": payload.RefreshToken,
		})
	}
}

var upgrader = gorilla.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, validate origin
	},
}

func wsHandler(hub *websocket.Hub, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantIDStr := c.Query("merchant_id")
		if merchantIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "merchant_id query parameter required"})
			return
		}

		merchantID, err := uuid.Parse(merchantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchant_id"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Error().Err(err).Msg("websocket upgrade failed")
			return
		}

		client := websocket.NewClient(hub, conn, merchantID)
		hub.Register(client)

		go client.WritePump()
		go client.ReadPump()
	}
}

const playgroundHTML = `<!DOCTYPE html>
<html>
<head>
  <title>Naseej GraphQL Playground</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/css/index.css" />
  <script src="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/js/middleware.js"></script>
</head>
<body>
  <div id="root"></div>
  <script>
    window.addEventListener('load', function() {
      GraphQLPlayground.init(document.getElementById('root'), {
        endpoint: '/graphql',
        settings: { 'request.credentials': 'include' }
      })
    })
  </script>
</body>
</html>`
