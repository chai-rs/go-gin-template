package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xanonydxck/simple-bookstore/api"
	"github.com/0xanonydxck/simple-bookstore/config"
	"github.com/0xanonydxck/simple-bookstore/docs"
	"github.com/0xanonydxck/simple-bookstore/infrastructure/auth"
	"github.com/0xanonydxck/simple-bookstore/infrastructure/db"
	_ "github.com/0xanonydxck/simple-bookstore/infrastructure/logger"
	eval "github.com/0xanonydxck/simple-bookstore/pkg/validator"
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
	app := gin.Default()

	// Setup middleware
	app.Use(logger.SetLogger())

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

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("🚨 failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("👋 server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("🚨 failed to shutdown server")
	}

	<-ctx.Done()
	log.Info().Msg("👋 server exited properly")
}

func setupSwagger(app *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "Simple Bookstore API"
	docs.SwaggerInfo.Description = "API for managing books in a bookstore"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", config.PORT)

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func setupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date_valid", eval.DateValid)
	}
}
