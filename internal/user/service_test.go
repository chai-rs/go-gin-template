package user_test

import (
	"context"
	"testing"

	"github.com/chai-rs/simple-bookstore/infrastructure/auth"
	errs "github.com/chai-rs/simple-bookstore/internal/error"
	"github.com/chai-rs/simple-bookstore/internal/model"
	"github.com/chai-rs/simple-bookstore/internal/user"
	"github.com/chai-rs/simple-bookstore/pkg/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestService_Login(t *testing.T) {
	type TestcaseIn struct {
		Email    string
		Password string
	}

	type Testcase struct {
		Name      string
		In        TestcaseIn
		WantError bool
	}

	testcases := []Testcase{
		{
			Name: "success",
			In: TestcaseIn{
				Email:    "one@example.com",
				Password: "password",
			},
		},
		{
			Name: "no-record",
			In: TestcaseIn{
				Email:    "invalid@example.com",
				Password: "password",
			},
			WantError: true,
		},
		{
			Name: "invalid-password",
			In: TestcaseIn{
				Email:    "one@example.com",
				Password: "invalid-password",
			},
			WantError: true,
		},
	}

	repo := user.NewMockRepository(t)
	repo.EXPECT().GetByEmail(mock.Anything, "one@example.com").Return(&model.User{
		ID:             uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Email:          "one@example.com",
		HashedPassword: "$2a$10$oiLJvjZFetwKPC5Gr9lBjuWuNdCYxorIsGJlSZtuhlnKmm4FxAoV6",
	}, nil)
	repo.EXPECT().GetByEmail(mock.Anything, "invalid@example.com").Return(nil, errs.FromGorm(gorm.ErrRecordNotFound))

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			tokenManager := auth.NewTokenManager()
			enforcer := auth.NewMockAuthEnforcer(t)
			memoryAuth := auth.NewMemoryAuth()

			svc := user.NewService(repo, memoryAuth, tokenManager, enforcer)
			accessToken, _, err := svc.Login(ctx, tc.In.Email, tc.In.Password)

			if tc.WantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				token, err := auth.VerifyToken(accessToken)
				assert.NoError(t, err)

				tokenProperties, err := auth.Extract(token)
				assert.NoError(t, err)

				_, err = memoryAuth.FetchAuth(ctx, tokenProperties.TokenUUID)
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_Register(t *testing.T) {
	type Testcase struct {
		Name      string
		In        *model.User
		WantError bool
	}

	testcases := []Testcase{
		{
			Name: "success",
			In: &model.User{
				Email:          "one@example.com",
				HashedPassword: crypto.MustHashPassword("password"),
			},
		},
		{
			Name: "duplicate",
			In: &model.User{
				Email:          "duplicate@example.com",
				HashedPassword: crypto.MustHashPassword("password"),
			},
			WantError: true,
		},
	}

	repo := user.NewMockRepository(t)
	repo.EXPECT().Create(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, user *model.User) error {
		switch user.Email {
		case "duplicate@example.com":
			return errs.FromGorm(gorm.ErrDuplicatedKey)
		default:
			return nil
		}
	})

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			enforcer := auth.NewMockAuthEnforcer(t)
			enforcer.EXPECT().AddPolicy(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

			tokenManager := auth.NewTokenManager()
			memoryAuth := auth.NewMemoryAuth()

			svc := user.NewService(repo, memoryAuth, tokenManager, enforcer)
			accessToken, _, err := svc.Register(ctx, tc.In)

			if tc.WantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				token, err := auth.VerifyToken(accessToken)
				assert.NoError(t, err)

				tokenProperties, err := auth.Extract(token)
				assert.NoError(t, err)

				userId, err := memoryAuth.FetchAuth(ctx, tokenProperties.TokenUUID)
				assert.NoError(t, err)
				assert.Equal(t, tc.In.ID.String(), userId)
			}
		})
	}
}

func TestService_Logout(t *testing.T) {
	ctx := context.Background()

	repo := user.NewMockRepository(t)
	repo.EXPECT().GetByEmail(mock.Anything, mock.Anything).Return(&model.User{
		ID:             uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Email:          "one@example.com",
		HashedPassword: "$2a$10$oiLJvjZFetwKPC5Gr9lBjuWuNdCYxorIsGJlSZtuhlnKmm4FxAoV6",
	}, nil)

	tokenManager := auth.NewTokenManager()
	memoryAuth := auth.NewMemoryAuth()

	enforcer := auth.NewMockAuthEnforcer(t)
	enforcer.EXPECT().AddPolicy(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	svc := user.NewService(repo, memoryAuth, tokenManager, enforcer)
	accessToken, _, err := svc.Login(ctx, "one@example.com", "password")
	assert.NoError(t, err)

	token, err := auth.VerifyToken(accessToken)
	assert.NoError(t, err)

	tokenProperties, err := auth.Extract(token)
	assert.NoError(t, err)

	assert.NoError(t, svc.Logout(ctx, tokenProperties))

	_, err = memoryAuth.FetchAuth(ctx, tokenProperties.TokenUUID)
	assert.Error(t, err)
}
