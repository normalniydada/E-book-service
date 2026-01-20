package handler

import (
	"E-book-service/internal/domain"
	"E-book-service/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Handler теперь зависит от интерфейса ServiceInterface
type Handler struct {
	svc service.ServiceInterface
}

// NewHandler принимает интерфейс, что позволяет передавать в него моки в тестах
func NewHandler(s service.ServiceInterface) *Handler {
	return &Handler{svc: s}
}

// Вспомогательная функция для получения ID пользователя из JWT контекста
func getUID(c echo.Context) uint {
	val := c.Get("user_id")
	if val == nil {
		return 0
	}
	return val.(uint)
}

// --- PUBLIC ---

func (h *Handler) Health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (h *Handler) Register(c echo.Context) error {
	var r struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.Bind(&r); err != nil {
		return err
	}
	if err := h.svc.Register(r.Email, r.Password, r.Name); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Login(c echo.Context) error {
	var r struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&r); err != nil {
		return err
	}
	token, err := h.svc.Login(r.Email, r.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

// --- PROTECTED ---

// Profile
func (h *Handler) GetMe(c echo.Context) error {
	u, err := h.svc.GetProfile(getUID(c))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, u)
}

func (h *Handler) UpdateProfile(c echo.Context) error {
	var u domain.User
	if err := c.Bind(&u); err != nil {
		return err
	}
	u.ID = getUID(c)
	if err := h.svc.UpdateProfile(&u); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, u)
}

// Books
func (h *Handler) ListBooks(c echo.Context) error {
	books, err := h.svc.GetAllBooks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, books)
}

func (h *Handler) CreateBook(c echo.Context) error {
	var b domain.Book
	if err := c.Bind(&b); err != nil {
		return err
	}
	if err := h.svc.CreateBook(&b); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, b)
}

func (h *Handler) GetBook(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	b, err := h.svc.GetBook(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}
	return c.JSON(http.StatusOK, b)
}

func (h *Handler) UpdateBook(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	var b domain.Book
	if err := c.Bind(&b); err != nil {
		return err
	}
	b.ID = uint(id)
	if err := h.svc.UpdateBook(&b); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, b)
}

func (h *Handler) DeleteBook(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	if err := h.svc.DeleteBook(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetBookContent(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	b, err := h.svc.GetBook(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}
	return c.JSON(http.StatusOK, map[string]string{"content": b.Content})
}

// Authors
func (h *Handler) ListAuthors(c echo.Context) error {
	authors, err := h.svc.GetAllAuthors()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, authors)
}

func (h *Handler) CreateAuthor(c echo.Context) error {
	var a domain.Author
	if err := c.Bind(&a); err != nil {
		return err
	}
	if err := h.svc.CreateAuthor(&a); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, a)
}

func (h *Handler) GetAuthor(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	a, err := h.svc.GetAuthor(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Author not found"})
	}
	return c.JSON(http.StatusOK, a)
}

func (h *Handler) UpdateAuthor(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	var a domain.Author
	if err := c.Bind(&a); err != nil {
		return err
	}
	a.ID = uint(id)
	if err := h.svc.UpdateAuthor(&a); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, a)
}

func (h *Handler) DeleteAuthor(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	if err := h.svc.DeleteAuthor(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetAuthorBooks(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	books, err := h.svc.GetBooksByAuthor(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, books)
}

// Reviews
func (h *Handler) ListReviews(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	reviews, err := h.svc.GetReviews(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, reviews)
}

func (h *Handler) AddReview(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	var r domain.Review
	if err := c.Bind(&r); err != nil {
		return err
	}
	r.BookID = uint(id)
	r.UserID = getUID(c)
	if err := h.svc.AddReview(&r); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, r)
}

func (h *Handler) DeleteReview(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	if err := h.svc.DeleteReview(uint(id), getUID(c)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// Shelf
func (h *Handler) GetShelf(c echo.Context) error {
	shelf, err := h.svc.GetShelf(getUID(c))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, shelf)
}

func (h *Handler) AddToShelf(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	var r struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&r); err != nil {
		return err
	}
	if err := h.svc.SetShelfStatus(getUID(c), uint(id), r.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) RemoveFromShelf(c echo.Context) error {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	id := uint(idInt)
	if err := h.svc.RemoveFromShelf(getUID(c), uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
