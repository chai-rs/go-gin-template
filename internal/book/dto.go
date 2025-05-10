package book

import (
	"time"

	"github.com/0xanonydxck/simple-bookstore/internal/model"
	"github.com/google/uuid"
)

type CreateBookDTO struct {
	Title       string    `json:"title" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	GenreCode   string    `json:"genre_code" binding:"required"`
	TagCodes    []string  `json:"tag_codes" binding:"required"`
	ReleaseDate time.Time `json:"release_date" binding:"required,date_valid" time_format:"2006-01-02"`
}

func (c *CreateBookDTO) ToBook() *model.Book {
	genre := &model.Genre{
		Code: c.GenreCode,
	}

	tags := make([]model.Tag, len(c.TagCodes))
	for i, tagCode := range c.TagCodes {
		tags[i] = model.Tag{
			Code: tagCode,
		}
	}

	return &model.Book{
		Title:       c.Title,
		Author:      c.Author,
		Genre:       genre,
		Tags:        tags,
		ReleaseDate: &c.ReleaseDate,
	}
}

type UpdateBookDTO struct {
	Title       string    `json:"title" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	GenreCode   string    `json:"genre_code" binding:"required"`
	TagCodes    []string  `json:"tag_codes" binding:"required"`
	ReleaseDate time.Time `json:"release_date" binding:"required,date_valid" time_format:"2006-01-02"`
}

func (u *UpdateBookDTO) ToBook() *model.Book {
	genre := &model.Genre{
		Code: u.GenreCode,
	}

	tags := make([]model.Tag, len(u.TagCodes))
	for i, tagCode := range u.TagCodes {
		tags[i] = model.Tag{
			Code: tagCode,
		}
	}

	return &model.Book{
		Title:       u.Title,
		Author:      u.Author,
		Genre:       genre,
		Tags:        tags,
		ReleaseDate: &u.ReleaseDate,
	}
}

type TagDTO struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type GenreDTO struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type BookDTO struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Genre       GenreDTO  `json:"genre"`
	Tags        []TagDTO  `json:"tags"`
	ReleaseDate time.Time `json:"release_date"`
}

func FromBook(book *model.Book) *BookDTO {
	tags := make([]TagDTO, len(book.Tags))
	for i, tag := range book.Tags {
		tags[i] = TagDTO{Code: tag.Code, Name: tag.Name}
	}

	return &BookDTO{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Genre:       GenreDTO{Code: book.Genre.Code, Name: book.Genre.Name},
		Tags:        tags,
		ReleaseDate: *book.ReleaseDate,
	}
}
