package book

import (
	"net/http"

	"github.com/0xanonydxck/simple-bookstore/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler represents the HTTP handler for book operations
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance
func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// CreateBook godoc
// @Summary Create a new book
// @Description Add a new book to the bookstore
// @Tags books
// @Accept json
// @Produce json
// @Param book body CreateBookDTO true "Book information"
// @Param Authorization header string true "Bearer token"
// @Success 201 {object} model.Book
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /books [post]
func (h *Handler) CreateBook(c *gin.Context) {
	var createBookDTO CreateBookDTO
	if err := c.ShouldBindJSON(&createBookDTO); err != nil {
		utils.ResponseErrorWithStatus(c, http.StatusBadRequest, "invalid request body")
		return
	}

	book := createBookDTO.ToBook()
	if err := h.service.Create(c.Request.Context(), book); err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseCreated(c, book)
}

// GetBooks godoc
// @Summary Get all books
// @Description Retrieve all books from the bookstore
// @Tags books
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} model.Book
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /books [get]
func (h *Handler) GetBooks(c *gin.Context) {
	books, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	result := make([]*BookDTO, len(books))
	for i, book := range books {
		result[i] = FromBook(&book)
	}

	utils.ResponseOk(c, result)
}

// GetBook godoc
// @Summary Get a book by ID
// @Description Retrieve a book by its UUID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} model.Book
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /books/{id} [get]
func (h *Handler) GetBook(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		utils.ResponseErrorWithStatus(c, http.StatusBadRequest, "invalid book id")
		return
	}

	book, err := h.service.GetByID(c.Request.Context(), parsedId)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	bookDTO := FromBook(book)

	utils.ResponseOk(c, bookDTO)
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update an existing book's information
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param book body UpdateBookDTO true "Updated book information"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} model.Book
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /books/{id} [put]
func (h *Handler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		utils.ResponseErrorWithStatus(c, http.StatusBadRequest, "invalid book id")
		return
	}

	var updateBookDTO UpdateBookDTO
	if err := c.ShouldBindJSON(&updateBookDTO); err != nil {
		utils.ResponseErrorWithStatus(c, http.StatusBadRequest, "invalid request body")
		return
	}

	book := updateBookDTO.ToBook()
	book.ID = parsedId
	if err := h.service.Update(c.Request.Context(), book); err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseOk(c, book)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Remove a book from the bookstore
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} nil
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /books/{id} [delete]
func (h *Handler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		utils.ResponseErrorWithStatus(c, http.StatusBadRequest, "invalid book id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), parsedId); err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseOk(c, nil)
}
