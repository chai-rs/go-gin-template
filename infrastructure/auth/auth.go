package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// obj
type AuthObject string

func (o AuthObject) String() string {
	return string(o)
}

const (
	Resource = AuthObject("resource")
)

// act
type AuthAction string

func (a AuthAction) String() string {
	return string(a)
}

const (
	Write = AuthAction("write")
	Read  = AuthAction("read")
)

type AccessProperties struct {
	TokenUUID string
	UserID    string
	Email     string
}

type TokenProperties struct {
	AccessToken        string
	RefreshToken       string
	AccessTokenUUID    string
	RefreshTokenUUID   string
	AccessTokenExpire  int64
	RefreshTokenExpire int64
}

type Auth interface {
	CreateAuth(ctx context.Context, userId string, properties *TokenProperties) error
	FetchAuth(ctx context.Context, userId string) (string, error)
	DeleteRefreshToken(ctx context.Context, userId string) error
	DeleteAccessToken(ctx context.Context, properties *AccessProperties) error
}

type RedisAuth struct {
	client *redis.Client
}

func NewRedisAuth(client *redis.Client) *RedisAuth {
	return &RedisAuth{client}
}

func (r *RedisAuth) CreateAuth(ctx context.Context, userId string, properties *TokenProperties) error {
	at := time.Unix(properties.AccessTokenExpire, 0)
	rt := time.Unix(properties.RefreshTokenExpire, 0)
	now := time.Now()

	atCreated, err := r.client.Set(ctx, properties.AccessTokenUUID, userId, at.Sub(now)).Result()
	if err != nil {
		return err
	}

	rtCreated, err := r.client.Set(ctx, properties.RefreshTokenUUID, userId, rt.Sub(now)).Result()
	if err != nil {
		return err
	}

	if atCreated == "0" || rtCreated == "0" {
		return fmt.Errorf("failed to create auth")
	}

	return nil
}

func (r *RedisAuth) FetchAuth(ctx context.Context, tokenUUID string) (string, error) {
	userId, err := r.client.Get(ctx, tokenUUID).Result()
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (r *RedisAuth) DeleteAccessToken(ctx context.Context, properties *AccessProperties) error {
	refreshUUID := ToRefreshUUID(properties.TokenUUID, properties.UserID)
	deletedAt, err := r.client.Del(ctx, properties.TokenUUID).Result()
	if err != nil {
		return err
	}

	deletedRt, err := r.client.Del(ctx, refreshUUID).Result()
	if err != nil {
		return err
	}

	if deletedAt != 1 || deletedRt != 1 {
		return fmt.Errorf("failed to delete access token")
	}

	return nil
}

func (r *RedisAuth) DeleteRefreshToken(ctx context.Context, refreshUUID string) error {
	deleted, err := r.client.Del(ctx, refreshUUID).Result()
	if err != nil {
		return err
	}

	if deleted == 0 {
		return fmt.Errorf("failed to delete refresh token")
	}

	return nil
}

type MemoryAuth struct {
	storage sync.Map
}

func NewMemoryAuth() *MemoryAuth {
	return &MemoryAuth{}
}

func (m *MemoryAuth) CreateAuth(ctx context.Context, userId string, properties *TokenProperties) error {
	m.storage.Store(properties.AccessTokenUUID, userId)
	m.storage.Store(properties.RefreshTokenUUID, userId)
	return nil
}

func (m *MemoryAuth) FetchAuth(ctx context.Context, tokenUUID string) (string, error) {
	userId, ok := m.storage.Load(tokenUUID)
	if !ok {
		return "", fmt.Errorf("failed to fetch auth")
	}
	return userId.(string), nil
}

func (m *MemoryAuth) DeleteAccessToken(ctx context.Context, properties *AccessProperties) error {
	m.storage.Delete(properties.TokenUUID)
	m.storage.Delete(ToRefreshUUID(properties.TokenUUID, properties.UserID))
	return nil
}

func (m *MemoryAuth) DeleteRefreshToken(ctx context.Context, refreshUUID string) error {
	m.storage.Delete(refreshUUID)
	return nil
}

func ToRefreshUUID(tokenUUID, userId string) string {
	return fmt.Sprintf("%s++%s", tokenUUID, userId)
}
