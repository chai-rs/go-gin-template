package model

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID  `gorm:"column:id;primaryKey"`
	Title       string     `gorm:"column:title"`
	Author      string     `gorm:"column:author"`
	GenreCode   string     `gorm:"column:genre_code;index"`
	ReleaseDate *time.Time `gorm:"column:release_date"`
	CreatedAt   *time.Time `gorm:"column:created_at"`

	Genre *Genre `gorm:"foreignKey:GenreCode;references:Code"`
	Tags  []Tag  `gorm:"many2many:book_tags;joinForeignKey:BookID;joinReferences:TagCode"`
}

func (b *Book) TableName() string {
	return "books"
}

type Genre struct {
	Code string `gorm:"column:code;primaryKey;unique"`
	Name string `gorm:"column:name"`
}

func (g *Genre) TableName() string {
	return "genres"
}

type Tag struct {
	Code string `gorm:"column:code;primaryKey;unique"`
	Name string `gorm:"column:name"`
}

func (t *Tag) TableName() string {
	return "tags"
}
