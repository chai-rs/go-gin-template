package api

import (
	"github.com/0xanonydxck/simple-bookstore/config"
	"github.com/0xanonydxck/simple-bookstore/infrastructure/auth"
	"github.com/0xanonydxck/simple-bookstore/infrastructure/db"
	"github.com/0xanonydxck/simple-bookstore/infrastructure/limiter"
	"github.com/0xanonydxck/simple-bookstore/internal/book"
	"github.com/0xanonydxck/simple-bookstore/internal/middleware"
	"github.com/0xanonydxck/simple-bookstore/internal/user"
	"github.com/gin-gonic/gin"
)

// BindRoutes registers all API routes to the given router
func BindRoutes(router *gin.Engine, enforcer auth.AuthEnforcer) {
	api := router.Group("/api")
	api.Use(middleware.RateLimitMiddleware(limiter.NewMemoryLimiter(config.LIMIT_RATE)))

	authorized := api.Group("", middleware.AuthMiddleware())
	unauthorized := api.Group("")

	bindBookRoutes(authorized, enforcer)
	bindUserRoutes(authorized, unauthorized, enforcer)
}

// bindBookRoutes registers all book-related routes to the API router group
func bindBookRoutes(api *gin.RouterGroup, enforcer auth.AuthEnforcer) {
	router := api.Group("/books")
	hdl := book.NewHandler(book.NewService(book.NewRepository(db.PostgreSQL())))

	router.POST("", middleware.Authorize(auth.Resource, auth.Write, enforcer), hdl.CreateBook)
	router.GET("", middleware.Authorize(auth.Resource, auth.Read, enforcer), hdl.GetBooks)
	router.GET("/:id", middleware.Authorize(auth.Resource, auth.Read, enforcer), hdl.GetBook)
	router.PUT("/:id", middleware.Authorize(auth.Resource, auth.Write, enforcer), hdl.UpdateBook)
	router.DELETE("/:id", middleware.Authorize(auth.Resource, auth.Write, enforcer), hdl.DeleteBook)
}

// bindUserRoutes registers all user-related routes to the API router group
func bindUserRoutes(authorized, unauthorized *gin.RouterGroup, enforcer auth.AuthEnforcer) {
	hdl := user.NewHandler(
		user.NewService(
			user.NewRepository(db.PostgreSQL()),
			auth.NewRedisAuth(db.Redis()),
			auth.NewTokenManager(),
			enforcer,
		),
	)

	{
		router := unauthorized.Group("/users")
		router.POST("/login", hdl.Login)
		router.POST("/register", hdl.Register)
		router.POST("/refresh", hdl.RefreshToken)
	}

	{
		router := authorized.Group("/users")
		router.POST("/logout", hdl.Logout)
	}
}
