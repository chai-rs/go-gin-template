package api

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// BindSwaggerRoutes registers the /reference endpoint for serving the API reference UI.
func BindSwaggerRoutes(app *gin.Engine) {
	router := app.Group("")
	bindSwaggerRoutes(router)
}

// bindSwaggerRoutes sets up the /reference route that serves the Swagger UI.
func bindSwaggerRoutes(router *gin.RouterGroup) {
	router.GET("/reference", func(ctx *gin.Context) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple Bookstore API",
			},
			DarkMode: true,
		})

		if err != nil {
			log.Error().Err(err).Msg("🚨 failed to generate API reference")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API reference"})
			return
		}

		ctx.Data(http.StatusOK, "text/html", []byte(htmlContent))
	})
}
