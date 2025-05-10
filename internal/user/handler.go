package user

import (
	"github.com/0xanonydxck/simple-bookstore/infrastructure/auth"
	"github.com/0xanonydxck/simple-bookstore/internal/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Login godoc
// @Summary Login user
// @Description Authenticate a user and return access and refresh tokens
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginRequestDTO true "Login credentials"
// @Success 200 {object} LoginResponseDTO
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, err)
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseOk(c, LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Register godoc
// @Summary Register new user
// @Description Register a new user and return access and refresh tokens
// @Tags users
// @Accept json
// @Produce json
// @Param request body RegisterRequestDTO true "User registration data"
// @Success 200 {object} RegisterResponseDTO
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, err)
		return
	}

	user, err := req.ToUser()
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	accessToken, refreshToken, err := h.service.Register(c.Request.Context(), user)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseOk(c, RegisterResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout a user and delete access and refresh tokens
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} nil
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	metadata, err := auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	if err := h.service.Logout(c.Request.Context(), metadata); err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseOk(c, nil)
}

// RefreshToken godoc
// @Summary Refresh user token
// @Description Refresh a user's access token
// @Tags users
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequestDTO true "Refresh token"
// @Success 200 {object} RefreshTokenResponseDTO
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, err)
		return
	}

	accessToken, refreshToken, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseOk(c, RefreshTokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
