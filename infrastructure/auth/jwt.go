package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/0xanonydxck/simple-bookstore/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager interface {
	CreateToken(userId, email string) (*TokenProperties, error)
	ExtractTokenMetadata(*http.Request) (*AccessProperties, error)
}

type tokenManager struct{}

func NewTokenManager() TokenManager {
	return &tokenManager{}
}

func (t *tokenManager) CreateToken(userId, email string) (*TokenProperties, error) {
	properties := new(TokenProperties)

	now := time.Now()
	properties.AccessTokenExpire = now.Add(time.Minute * 30).Unix()
	properties.AccessTokenUUID = uuid.New().String()
	properties.RefreshTokenExpire = now.Add(time.Hour * 24 * 7).Unix()
	properties.RefreshTokenUUID = ToRefreshUUID(properties.RefreshTokenUUID, userId)

	// Create access token
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = properties.AccessTokenUUID
	atClaims["user_id"] = userId
	atClaims["email"] = email
	atClaims["exp"] = properties.AccessTokenExpire
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	properties.AccessToken, err = at.SignedString([]byte(config.ACCESS_SECRET))
	if err != nil {
		return nil, err
	}

	// Create refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = properties.RefreshTokenUUID
	rtClaims["user_id"] = userId
	rtClaims["email"] = email
	rtClaims["exp"] = properties.RefreshTokenExpire
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	properties.RefreshToken, err = rt.SignedString([]byte(config.REFRESH_SECRET))
	if err != nil {
		return nil, err
	}

	return properties, nil
}

func (t *tokenManager) ExtractTokenMetadata(r *http.Request) (*AccessProperties, error) {
	return ExtractTokenMetadata(r)
}

func TokenValid(r *http.Request) error {
	token, err := VerifyAuthorizationHeader(r)
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func VerifyAuthorizationHeader(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	return VerifyToken(tokenString)
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}

		return []byte(config.ACCESS_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func Extract(token *jwt.Token) (*AccessProperties, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	accessUUID, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid access uuid")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid email")
	}

	properties := &AccessProperties{
		TokenUUID: accessUUID,
		UserID:    userID,
		Email:     email,
	}

	return properties, nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessProperties, error) {
	token, err := VerifyAuthorizationHeader(r)
	if err != nil {
		return nil, err
	}

	acc, err := Extract(token)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
