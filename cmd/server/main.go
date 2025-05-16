package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/chai-rs/simple-bookstore/api"
	"github.com/chai-rs/simple-bookstore/config"
	"github.com/chai-rs/simple-bookstore/docs"
	"github.com/chai-rs/simple-bookstore/infrastructure/auth"
	"github.com/chai-rs/simple-bookstore/infrastructure/db"
	_ "github.com/chai-rs/simple-bookstore/infrastructure/logger"
	eval "github.com/chai-rs/simple-bookstore/pkg/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	config.Init()
	db.PostgreSQLConnect(
		config.POSTGRES_HOST,
		config.POSTGRES_PORT,
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.POSTGRES_DB,
	)

	setupValidator()
}

func main() {
	// Setup gin engine
	app := setupGin()

	// Setup middleware
	app.Use(logger.SetLogger())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(config.CORS_ALLOWED_ORIGINS, ","),
		AllowMethods:     strings.Split(config.CORS_ALLOWED_METHODS, ","),
		AllowHeaders:     strings.Split(config.CORS_ALLOWED_HEADERS, ","),
		ExposeHeaders:    strings.Split(config.CORS_EXPOSED_HEADERS, ","),
		MaxAge:           time.Duration(config.CORS_MAX_AGE),
		AllowCredentials: true,
	}))

	// Setup enforcer
	enforcer := auth.NewAuthEnforcer(auth.GormAdapter(db.PostgreSQL()))

	// Bind routes
	api.BindRoutes(app, enforcer)

	// Setup swagger
	setupSwagger(app)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.PORT),
		Handler: app.Handler(),
	}

	// Start server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("🚨 failed to start server")
		}
	}()

	// Wait for signal to stop server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("👋 server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("🚨 failed to shutdown server")
	}

	// Wait for server to exit
	<-ctx.Done()
	log.Info().Msg("👋 server exited properly")
}

// Setup Gin Engine
func setupGin() *gin.Engine {
	if config.MODE == config.ProductionMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	return gin.Default()
}

// Setup swagger
func setupSwagger(app *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "Simple Bookstore API"
	docs.SwaggerInfo.Description = "API for managing books in a bookstore"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", config.PORT)

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Setup validator
func setupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date_valid", eval.DateValid)
	}
}
