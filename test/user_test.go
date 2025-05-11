package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xanonydxck/simple-bookstore/infrastructure/auth"
	"github.com/0xanonydxck/simple-bookstore/internal/user"
	"github.com/0xanonydxck/simple-bookstore/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TestUserSuite runs the user suite.
func TestUserSuite(t *testing.T) {
	suite.Run(t, new(TestUserRegister_Successful))
	suite.Run(t, new(TestUserLogin_Successful))
}

// UserSuite runs the user suite.
type UserSuite struct {
	BaseSuite
	router  *gin.Engine
	service user.Service
}

// BeforeTest runs before each test.
func (s *UserSuite) BeforeTest(suiteName, testName string) {
	repo := user.NewRepository(s.db)
	enforcer := auth.NewAuthEnforcer(auth.GormAdapter(s.db), &auth.AuthEnforcerOpts{
		ModelPath: "../auth_model.conf",
	})
	tokenManager := auth.NewTokenManager()
	memoryAuth := auth.NewMemoryAuth()
	s.service = user.NewService(repo, memoryAuth, tokenManager, enforcer)
	hdl := user.NewHandler(s.service)

	s.router = gin.Default()
	s.router.POST("/register", hdl.Register)
	s.router.POST("/login", hdl.Login)
}

// TestUserRegister_Successful tests the user register successful.
type TestUserRegister_Successful struct{ UserSuite }

// TestUserRegister_Successful tests the user register successful.
func (s *TestUserRegister_Successful) TestUserRegister_Successful() {
	body := user.RegisterRequestDTO{
		Email:    "test@example.com",
		Password: "password",
	}

	jsonBody, err := json.Marshal(body)
	assert.NoError(s.T(), err)

	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonBody))
	assert.NoError(s.T(), err)

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
}

// TestUserLogin_Successful tests the user login successful.
type TestUserLogin_Successful struct{ UserSuite }

// BeforeTest runs before each test.
func (s *TestUserLogin_Successful) BeforeTest(suiteName, testName string) {
	s.UserSuite.BeforeTest(suiteName, testName)

	body := user.RegisterRequestDTO{
		Email:    "test@example.com",
		Password: "password",
	}

	user, err := body.ToUser()
	assert.NoError(s.T(), err)

	_, _, err = s.service.Register(context.Background(), user)
	assert.NoError(s.T(), err)
}

// TestUserLogin_Successful tests the user login successful.
func (s *TestUserLogin_Successful) TestUserLogin_Successful() {
	body := user.LoginRequestDTO{
		Email:    "test@example.com",
		Password: "password",
	}

	jsonBody, err := json.Marshal(body)
	assert.NoError(s.T(), err)

	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	assert.NoError(s.T(), err)

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	jsonBodyResp, err := io.ReadAll(w.Body)
	assert.NoError(s.T(), err)

	var resp utils.Response
	err = json.Unmarshal(jsonBodyResp, &resp)
	assert.NoError(s.T(), err)

	resultBuf, err := json.Marshal(resp.Result)
	assert.NoError(s.T(), err)

	var result user.LoginResponseDTO
	err = json.Unmarshal(resultBuf, &result)
	assert.NoError(s.T(), err)

	_, err = auth.VerifyToken(result.AccessToken)
	assert.NoError(s.T(), err)
}
