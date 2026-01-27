package handler

import (
	"E-book-service/internal/domain"
	"E-book-service/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ShelfStatusRequest struct {
	Status string `json:"status"`
}

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

// Health godoc
// @Summary Проверка работоспособности
// @Tags System
// @Success 200 {string} string "OK"
// @Router /health [get]
func (h *Handler) Health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// Register godoc
// @Summary Регистрация пользователя
// @Tags Auth
// @Accept json
// @Param body body RegisterRequest true "Данные регистрации"
// @Success 201 "Created"
// @Router /register [post]
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

// Login godoc
// @Summary Авторизация
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Данные логина"
// @Success 200 {object} map[string]string "token"
// @Router /login [post]
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

// @Summary Получить свой профиль
// @Tags Profile
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} domain.User
// @Router /me [get]
func (h *Handler) GetMe(c echo.Context) error {
	u, err := h.svc.GetProfile(getUID(c))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, u)
}

// @Summary Обновить профиль
// @Tags Profile
// @Security ApiKeyAuth
// @Accept json
// @Param user body domain.User true "Данные профиля"
// @Success 200 {object} domain.User
// @Router /me [put]
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

// @Summary Список всех книг
// @Tags Books
// @Produce json
// @Success 200 {array} domain.Book
// @Router /books [get]
func (h *Handler) ListBooks(c echo.Context) error {
	books, err := h.svc.GetAllBooks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, books)
}

// @Summary Создать книгу
// @Tags Books
// @Security ApiKeyAuth
// @Accept json
// @Param book body domain.Book true "Данные книги"
// @Success 201 {object} domain.Book
// @Router /books [post]
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

// @Summary Получить книгу по ID
// @Tags Books
// @Param id path int true "ID книги"
// @Produce json
// @Success 200 {object} domain.Book
// @Router /books/{id} [get]
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

// @Summary Обновить книгу
// @Tags Books
// @Security ApiKeyAuth
// @Param id path int true "ID книги"
// @Accept json
// @Param book body domain.Book true "Новые данные"
// @Success 200 {object} domain.Book
// @Router /books/{id} [put]
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

// @Summary Удалить книгу
// @Tags Books
// @Security ApiKeyAuth
// @Param id path int true "ID книги"
// @Success 204 "No Content"
// @Router /books/{id} [delete]
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

// @Summary Текст книги
// @Tags Books
// @Param id path int true "ID книги"
// @Produce json
// @Success 200 {object} map[string]string "content"
// @Router /books/{id}/content [get]
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

// @Summary Список авторов
// @Tags Authors
// @Produce json
// @Success 200 {array} domain.Author
// @Router /authors [get]
func (h *Handler) ListAuthors(c echo.Context) error {
	authors, err := h.svc.GetAllAuthors()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, authors)
}

// @Summary Создать автора
// @Tags Authors
// @Security ApiKeyAuth
// @Accept json
// @Param author body domain.Author true "Данные автора"
// @Success 201 {object} domain.Author
// @Router /authors [post]
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

// @Summary Инфо об авторе
// @Tags Authors
// @Param id path int true "ID автора"
// @Produce json
// @Success 200 {object} domain.Author
// @Router /authors/{id} [get]
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

// @Summary Обновить автора
// @Tags Authors
// @Security ApiKeyAuth
// @Param id path int true "ID автора"
// @Accept json
// @Param author body domain.Author true "Данные"
// @Success 200 {object} domain.Author
// @Router /authors/{id} [put]
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

// @Summary Удалить автора
// @Tags Authors
// @Security ApiKeyAuth
// @Param id path int true "ID автора"
// @Success 204 "No Content"
// @Router /authors/{id} [delete]
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

// @Summary Книги автора
// @Tags Authors
// @Param id path int true "ID автора"
// @Produce json
// @Success 200 {array} domain.Book
// @Router /authors/{id}/books [get]
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

// @Summary Список отзывов к книге
// @Tags Reviews
// @Param id path int true "ID книги"
// @Produce json
// @Success 200 {array} domain.Review
// @Router /books/{id}/reviews [get]
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

// @Summary Добавить отзыв
// @Tags Reviews
// @Security ApiKeyAuth
// @Param id path int true "ID книги"
// @Accept json
// @Param review body domain.Review true "Отзыв"
// @Success 201 {object} domain.Review
// @Router /books/{id}/reviews [post]
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

// @Summary Удалить отзыв
// @Tags Reviews
// @Security ApiKeyAuth
// @Param id path int true "ID отзыва"
// @Success 204 "No Content"
// @Router /reviews/{id} [delete]
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

// @Summary Моя полка
// @Tags Shelf
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} domain.Book
// @Router /shelf [get]
func (h *Handler) GetShelf(c echo.Context) error {
	shelf, err := h.svc.GetShelf(getUID(c))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, shelf)
}

// AddToShelf godoc
// @Summary Добавить на полку
// @Tags Shelf
// @Security ApiKeyAuth
// @Param id path int true "Book ID"
// @Accept json
// @Param body body ShelfStatusRequest true "Статус"
// @Router /shelf/{id} [post]
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

// @Summary Удалить с полки
// @Tags Shelf
// @Security ApiKeyAuth
// @Param id path int true "ID книги"
// @Success 204 "No Content"
// @Router /shelf/{id} [delete]
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
