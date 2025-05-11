package book_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/0xanonydxck/simple-bookstore/internal/book"
	errs "github.com/0xanonydxck/simple-bookstore/internal/error"
	"github.com/0xanonydxck/simple-bookstore/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"go.openly.dev/pointy"
)

func TestService_Create(t *testing.T) {
	type Testcase struct {
		Name      string
		In        *model.Book
		WantError bool
	}

	testcases := []Testcase{
		{
			Name: "success",
			In: &model.Book{
				Title:       "Book 1",
				Author:      "Author 1",
				Genre:       &model.Genre{Code: "genre1", Name: "Genre 1"},
				Tags:        []model.Tag{{Code: "tag1", Name: "Tag 1"}},
				ReleaseDate: pointy.Pointer(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			Name: "invalid-title",
			In: &model.Book{
				Title:       "",
				Author:      "Author 1",
				Genre:       &model.Genre{Code: "genre1", Name: "Genre 1"},
				Tags:        []model.Tag{{Code: "tag1", Name: "Tag 1"}},
				ReleaseDate: pointy.Pointer(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			WantError: true,
		},
		{
			Name:      "empty",
			In:        nil,
			WantError: true,
		},
		{
			Name: "duplicated",
			In: &model.Book{
				Title:       "duplicated",
				Author:      "Author 1",
				Genre:       &model.Genre{Code: "genre1", Name: "Genre 1"},
				Tags:        []model.Tag{{Code: "tag1", Name: "Tag 1"}},
				ReleaseDate: pointy.Pointer(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			WantError: true,
		},
	}

	repo := book.NewMockRepository(t)
	repo.EXPECT().
		Create(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, book *model.Book) error {
			if book == nil {
				return errs.New(http.StatusBadRequest, fmt.Errorf("book is nil"))
			}

			if book.Title == "duplicated" {
				return errs.New(http.StatusBadRequest, fmt.Errorf("book already exists"))
			}

			v := validator.New()
			if err := v.Var(book.Title, "required"); err != nil {
				return errs.New(http.StatusBadRequest, err)
			}

			if err := v.Var(book.Author, "required"); err != nil {
				return errs.New(http.StatusBadRequest, err)
			}

			if err := v.Var(book.Genre, "required"); err != nil {
				return errs.New(http.StatusBadRequest, err)
			}

			if err := v.Var(book.Tags, "required"); err != nil {
				return errs.New(http.StatusBadRequest, err)
			}

			if err := v.Var(book.ReleaseDate, "required"); err != nil {
				return errs.New(http.StatusBadRequest, err)
			}

			return nil
		}).Maybe()

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			svc := book.NewService(repo)
			err := svc.Create(ctx, tc.In)

			if tc.WantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	type Testcase struct {
		Name      string
		In        *model.Book
		WantError bool
	}

	testcases := []Testcase{
		{
			Name: "success",
			In: &model.Book{
				ID:          uuid.New(),
				Title:       "Book 1",
				Author:      "Author 1",
				Genre:       &model.Genre{Code: "genre1", Name: "Genre 1"},
				Tags:        []model.Tag{{Code: "tag1", Name: "Tag 1"}},
				ReleaseDate: pointy.Pointer(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			Name: "invalid-id",
			In: &model.Book{
				Title:       "Book 1",
				Author:      "Author 1",
				Genre:       &model.Genre{Code: "genre1", Name: "Genre 1"},
				Tags:        []model.Tag{{Code: "tag1", Name: "Tag 1"}},
				ReleaseDate: pointy.Pointer(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			WantError: true,
		},
		{
			Name: "invalid-title",
			In: &model.Book{
				ID:          uuid.New(),
				Title:       "",
				Author:      "Author 1",
				Genre:       &model.Genre{Code: "genre1", Name: "Genre 1"},
				Tags:        []model.Tag{{Code: "tag1", Name: "Tag 1"}},
				ReleaseDate: pointy.Pointer(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			WantError: true,
		},
		{
			Name:      "empty",
			In:        nil,
			WantError: true,
		},
	}

	repo := book.NewMockRepository(t)
	repo.EXPECT().
		Update(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, book *model.Book) error {
			if book == nil {
				return errs.New(http.StatusBadRequest, fmt.Errorf("book is nil"))
			}

			if book.ID == uuid.Nil {
				return errs.New(http.StatusBadRequest, fmt.Errorf("book id is nil"))
			}

			if book.Title == "duplicated" {
				return errs.New(http.StatusBadRequest, fmt.Errorf("book already exists"))
			}

			v := validator.New()
			if err := v.Var(book.Title, "required"); err != nil {
				return errs.New(http.StatusBadRequest, err)
			}

			if err := v.Var(book.Author, "required"); err != nil {
				return errs.New(http.StatusBadRequest, err)
			}

			if err := v.Var(book.Genre, "required"); err != nil {
				return errs.New(http.StatusBadRequest, err)
			}

			if err := v.Var(book.Tags, "required"); err != nil {
				return errs.New(http.StatusBadRequest, err)
			}

			if err := v.Var(book.ReleaseDate, "required"); err != nil {
				return errs.New(http.StatusBadRequest, err)
			}

			return nil
		}).Maybe()

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			svc := book.NewService(repo)
			err := svc.Update(ctx, tc.In)

			if tc.WantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetAll(t *testing.T) {
	type Testcase struct {
		Name      string
		WantError bool
	}

	testcases := []Testcase{
		{
			Name: "success",
		},
	}

	data := []model.Book{
		{
			ID:          uuid.New(),
			Title:       "Book 1",
			Author:      "Author 1",
			Genre:       &model.Genre{Code: "genre1", Name: "Genre 1"},
			Tags:        []model.Tag{{Code: "tag1", Name: "Tag 1"}},
			ReleaseDate: pointy.Pointer(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	repo := book.NewMockRepository(t)
	repo.EXPECT().
		GetAll(mock.Anything).
		RunAndReturn(func(ctx context.Context) ([]model.Book, error) {
			return data, nil
		}).Maybe()

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			svc := book.NewService(repo)
			books, err := svc.GetAll(ctx)

			if tc.WantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, data, books)
			}
		})
	}
}

func TestService_GetByID(t *testing.T) {
	data := []model.Book{
		{
			ID:          uuid.New(),
			Title:       "Book 1",
			Author:      "Author 1",
			Genre:       &model.Genre{Code: "genre1", Name: "Genre 1"},
			Tags:        []model.Tag{{Code: "tag1", Name: "Tag 1"}},
			ReleaseDate: pointy.Pointer(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	type Testcase struct {
		Name      string
		In        uuid.UUID
		WantError bool
	}

	testcases := []Testcase{
		{
			Name: "success",
			In:   data[0].ID,
		},
		{
			Name:      "not-found",
			In:        uuid.New(),
			WantError: true,
		},
		{
			Name:      "empty",
			In:        uuid.Nil,
			WantError: true,
		},
	}

	repo := book.NewMockRepository(t)
	repo.EXPECT().
		GetByID(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, id uuid.UUID) (*model.Book, error) {
			if id == uuid.Nil {
				return nil, errs.New(http.StatusBadRequest, fmt.Errorf("book id is nil"))
			}

			var result *model.Book
			for _, book := range data {
				if book.ID == id {
					result = &book
					break
				}
			}

			if result == nil {
				return nil, errs.New(http.StatusNotFound, fmt.Errorf("book not found"))
			}

			return result, nil
		}).Maybe()

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			svc := book.NewService(repo)
			book, err := svc.GetByID(ctx, tc.In)

			if tc.WantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, data[0], *book)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	type Testcase struct {
		Name      string
		In        uuid.UUID
		WantError bool
	}

	testcases := []Testcase{
		{
			Name: "success",
			In:   uuid.New(),
		},
		{
			Name:      "empty",
			In:        uuid.Nil,
			WantError: true,
		},
	}

	repo := book.NewMockRepository(t)
	repo.EXPECT().
		Delete(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, id uuid.UUID) error {
			if id == uuid.Nil {
				return errs.New(http.StatusBadRequest, fmt.Errorf("book id is nil"))
			}

			return nil
		}).Maybe()

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			svc := book.NewService(repo)
			err := svc.Delete(ctx, tc.In)

			if tc.WantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
