package handler

import (
	"E-book-service/internal/domain"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCK SERVICE ---

type MockService struct {
	mock.Mock
}

func (m *MockService) Register(email, pass, name string) error {
	return m.Called(email, pass, name).Error(0)
}
func (m *MockService) Login(email, pass string) (string, error) {
	args := m.Called(email, pass)
	return args.String(0), args.Error(1)
}
func (m *MockService) GetProfile(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockService) UpdateProfile(u *domain.User) error { return m.Called(u).Error(0) }
func (m *MockService) CreateBook(b *domain.Book) error    { return m.Called(b).Error(0) }
func (m *MockService) GetAllBooks() ([]domain.Book, error) {
	args := m.Called()
	return args.Get(0).([]domain.Book), args.Error(1)
}
func (m *MockService) GetBook(id uint) (*domain.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Book), args.Error(1)
}
func (m *MockService) UpdateBook(b *domain.Book) error { return m.Called(b).Error(0) }
func (m *MockService) DeleteBook(id uint) error        { return m.Called(id).Error(0) }
func (m *MockService) GetBooksByAuthor(aID uint) ([]domain.Book, error) {
	args := m.Called(aID)
	return args.Get(0).([]domain.Book), args.Error(1)
}
func (m *MockService) CreateAuthor(a *domain.Author) error { return m.Called(a).Error(0) }
func (m *MockService) GetAllAuthors() ([]domain.Author, error) {
	args := m.Called()
	return args.Get(0).([]domain.Author), args.Error(1)
}
func (m *MockService) GetAuthor(id uint) (*domain.Author, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Author), args.Error(1)
}
func (m *MockService) UpdateAuthor(a *domain.Author) error { return m.Called(a).Error(0) }
func (m *MockService) DeleteAuthor(id uint) error          { return m.Called(id).Error(0) }
func (m *MockService) AddReview(re *domain.Review) error   { return m.Called(re).Error(0) }
func (m *MockService) GetReviews(bID uint) ([]domain.Review, error) {
	args := m.Called(bID)
	return args.Get(0).([]domain.Review), args.Error(1)
}
func (m *MockService) DeleteReview(id, uID uint) error { return m.Called(id, uID).Error(0) }
func (m *MockService) SetShelfStatus(uID, bID uint, status string) error {
	return m.Called(uID, bID, status).Error(0)
}
func (m *MockService) GetShelf(uID uint) ([]domain.Shelf, error) {
	args := m.Called(uID)
	return args.Get(0).([]domain.Shelf), args.Error(1)
}
func (m *MockService) RemoveFromShelf(uID, bID uint) error { return m.Called(uID, bID).Error(0) }

// --- TESTS ---

func TestHandler_All(t *testing.T) {
	e := echo.New()
	ms := new(MockService)
	h := NewHandler(ms)

	t.Run("Health", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		assert.NoError(t, h.Health(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Auth_Register_BindErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/reg", strings.NewReader("{invalid}"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, httptest.NewRecorder())
		assert.Error(t, h.Register(c))
	})

	t.Run("Auth_Register_SvcErr", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"email": "e", "password": "p"})
		req := httptest.NewRequest(http.MethodPost, "/reg", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ms.On("Register", "e", "p", "").Return(errors.New("fail")).Once()
		assert.NoError(t, h.Register(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Auth_Login_Success", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"email": "e", "password": "p"})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ms.On("Login", "e", "p").Return("token", nil).Once()
		assert.NoError(t, h.Login(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Auth_Login_BindErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("!"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, httptest.NewRecorder())
		assert.Error(t, h.Login(c))
	})

	t.Run("Profile_GetMe_NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", uint(1))
		ms.On("GetProfile", uint(1)).Return(nil, errors.New("err")).Once()
		assert.NoError(t, h.GetMe(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Profile_GetMe_NoUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec) // No UID set
		ms.On("GetProfile", uint(0)).Return(nil, errors.New("err")).Once()
		assert.NoError(t, h.GetMe(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Profile_Update_SvcErr", func(t *testing.T) {
		body, _ := json.Marshal(domain.User{Name: "N"})
		req := httptest.NewRequest(http.MethodPut, "/me", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", uint(1))
		ms.On("UpdateProfile", mock.Anything).Return(errors.New("err")).Once()
		assert.NoError(t, h.UpdateProfile(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Profile_Update_BindErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/me", strings.NewReader("?"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, httptest.NewRecorder())
		assert.Error(t, h.UpdateProfile(c))
	})

	t.Run("Books_List_Err", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ms.On("GetAllBooks").Return([]domain.Book{}, errors.New("err")).Once()
		assert.NoError(t, h.ListBooks(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Books_Create_SvcErr", func(t *testing.T) {
		body, _ := json.Marshal(domain.Book{Title: "T"})
		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ms.On("CreateBook", mock.Anything).Return(errors.New("err")).Once()
		assert.NoError(t, h.CreateBook(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Books_Create_BindErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader("?"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, httptest.NewRecorder())
		assert.Error(t, h.CreateBook(c))
	})

	t.Run("Books_Get_NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		ms.On("GetBook", uint(1)).Return(nil, errors.New("err")).Once()
		assert.NoError(t, h.GetBook(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Books_Update_SvcErr", func(t *testing.T) {
		body, _ := json.Marshal(domain.Book{Title: "U"})
		req := httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		ms.On("UpdateBook", mock.Anything).Return(errors.New("err")).Once()
		assert.NoError(t, h.UpdateBook(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Books_Update_BindErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/books/1", strings.NewReader("?"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, httptest.NewRecorder())
		assert.Error(t, h.UpdateBook(c))
	})

	t.Run("Books_Delete_SvcErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/books/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		ms.On("DeleteBook", uint(1)).Return(errors.New("err")).Once()
		assert.NoError(t, h.DeleteBook(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Books_GetContent_NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books/1/content", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		ms.On("GetBook", uint(1)).Return(nil, errors.New("err")).Once()
		assert.NoError(t, h.GetBookContent(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Authors_List_Err", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/authors", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ms.On("GetAllAuthors").Return([]domain.Author{}, errors.New("err")).Once()
		assert.NoError(t, h.ListAuthors(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Authors_Create_SvcErr", func(t *testing.T) {
		body, _ := json.Marshal(domain.Author{Name: "A"})
		req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ms.On("CreateAuthor", mock.Anything).Return(errors.New("err")).Once()
		assert.NoError(t, h.CreateAuthor(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Authors_Create_BindErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/authors", strings.NewReader("?"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, httptest.NewRecorder())
		assert.Error(t, h.CreateAuthor(c))
	})

	t.Run("Authors_Get_NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/authors/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		ms.On("GetAuthor", uint(1)).Return(nil, errors.New("err")).Once()
		assert.NoError(t, h.GetAuthor(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Authors_Update_SvcErr", func(t *testing.T) {
		body, _ := json.Marshal(domain.Author{Name: "U"})
		req := httptest.NewRequest(http.MethodPut, "/authors/1", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		ms.On("UpdateAuthor", mock.Anything).Return(errors.New("err")).Once()
		assert.NoError(t, h.UpdateAuthor(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Authors_Update_BindErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/authors/1", strings.NewReader("?"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, httptest.NewRecorder())
		assert.Error(t, h.UpdateAuthor(c))
	})

	t.Run("Authors_Delete_SvcErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/authors/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		ms.On("DeleteAuthor", uint(1)).Return(errors.New("err")).Once()
		assert.NoError(t, h.DeleteAuthor(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Authors_GetBooks_Err", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/authors/1/books", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		ms.On("GetBooksByAuthor", uint(1)).Return([]domain.Book{}, errors.New("err")).Once()
		assert.NoError(t, h.GetAuthorBooks(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Reviews_List_Err", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reviews/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		ms.On("GetReviews", uint(1)).Return([]domain.Review{}, errors.New("err")).Once()
		assert.NoError(t, h.ListReviews(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Reviews_Add_SvcErr", func(t *testing.T) {
		body, _ := json.Marshal(domain.Review{Comment: "C"})
		req := httptest.NewRequest(http.MethodPost, "/reviews/1", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("user_id", uint(1))
		ms.On("AddReview", mock.Anything).Return(errors.New("err")).Once()
		assert.NoError(t, h.AddReview(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Reviews_Add_BindErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/reviews/1", strings.NewReader("?"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, httptest.NewRecorder())
		assert.Error(t, h.AddReview(c))
	})

	t.Run("Reviews_Delete_SvcErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/reviews/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("user_id", uint(1))
		ms.On("DeleteReview", uint(1), uint(1)).Return(errors.New("err")).Once()
		assert.NoError(t, h.DeleteReview(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Shelf_Get_Err", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shelf", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", uint(1))
		ms.On("GetShelf", uint(1)).Return([]domain.Shelf{}, errors.New("err")).Once()
		assert.NoError(t, h.GetShelf(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Shelf_Add_BindErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/shelf/1", strings.NewReader("?"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, httptest.NewRecorder())
		assert.Error(t, h.AddToShelf(c))
	})

	t.Run("Shelf_Add_SvcErr", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"status": "s"})
		req := httptest.NewRequest(http.MethodPost, "/shelf/1", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("user_id", uint(1))
		ms.On("SetShelfStatus", uint(1), uint(1), "s").Return(errors.New("err")).Once()
		assert.NoError(t, h.AddToShelf(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Shelf_Remove_SvcErr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/shelf/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("user_id", uint(1))
		ms.On("RemoveFromShelf", uint(1), uint(1)).Return(errors.New("err")).Once()
		assert.NoError(t, h.RemoveFromShelf(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
