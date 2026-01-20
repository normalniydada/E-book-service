package service

import (
	"E-book-service/internal/domain"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(u *domain.User) error { return m.Called(u).Error(0) }
func (m *MockRepository) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockRepository) GetUserByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockRepository) UpdateUser(u *domain.User) error { return m.Called(u).Error(0) }

func (m *MockRepository) CreateBook(b *domain.Book) error { return m.Called(b).Error(0) }
func (m *MockRepository) GetBooks() ([]domain.Book, error) {
	args := m.Called()
	return args.Get(0).([]domain.Book), args.Error(1)
}
func (m *MockRepository) GetBookByID(id uint) (*domain.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Book), args.Error(1)
}
func (m *MockRepository) UpdateBook(b *domain.Book) error { return m.Called(b).Error(0) }
func (m *MockRepository) DeleteBook(id uint) error        { return m.Called(id).Error(0) }
func (m *MockRepository) GetBooksByAuthor(aID uint) ([]domain.Book, error) {
	args := m.Called(aID)
	return args.Get(0).([]domain.Book), args.Error(1)
}

func (m *MockRepository) CreateAuthor(a *domain.Author) error { return m.Called(a).Error(0) }
func (m *MockRepository) GetAuthors() ([]domain.Author, error) {
	args := m.Called()
	return args.Get(0).([]domain.Author), args.Error(1)
}
func (m *MockRepository) GetAuthorByID(id uint) (*domain.Author, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Author), args.Error(1)
}
func (m *MockRepository) UpdateAuthor(a *domain.Author) error { return m.Called(a).Error(0) }
func (m *MockRepository) DeleteAuthor(id uint) error          { return m.Called(id).Error(0) }

func (m *MockRepository) CreateReview(re *domain.Review) error { return m.Called(re).Error(0) }
func (m *MockRepository) GetReviewsByBook(bID uint) ([]domain.Review, error) {
	args := m.Called(bID)
	return args.Get(0).([]domain.Review), args.Error(1)
}
func (m *MockRepository) DeleteReview(id, uID uint) error { return m.Called(id, uID).Error(0) }

func (m *MockRepository) AddToShelf(s *domain.Shelf) error { return m.Called(s).Error(0) }
func (m *MockRepository) GetShelf(uID uint) ([]domain.Shelf, error) {
	args := m.Called(uID)
	return args.Get(0).([]domain.Shelf), args.Error(1)
}
func (m *MockRepository) RemoveFromShelf(uID, bID uint) error { return m.Called(uID, bID).Error(0) }

// --- ТЕСТЫ ---

func TestAuthAndProfile(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, "test-key")

	t.Run("Register", func(t *testing.T) {
		mockRepo.On("CreateUser", mock.Anything).Return(nil).Once()
		err := svc.Register("test@mail.com", "pass", "User")
		assert.NoError(t, err)
	})

	t.Run("Login_Success", func(t *testing.T) {
		hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
		mockRepo.On("GetUserByEmail", "test@mail.com").Return(&domain.User{Email: "test@mail.com", Password: string(hash)}, nil).Once()
		token, err := svc.Login("test@mail.com", "pass")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Login_Fail", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", "fail@mail.com").Return(nil, errors.New("not found")).Once()
		_, err := svc.Login("fail@mail.com", "any")
		assert.Error(t, err)
	})

	t.Run("ProfileOperations", func(t *testing.T) {
		user := &domain.User{Name: "Name"}
		mockRepo.On("GetUserByID", uint(1)).Return(user, nil).Once()
		res, _ := svc.GetProfile(1)
		assert.Equal(t, "Name", res.Name)

		mockRepo.On("UpdateUser", user).Return(nil).Once()
		err := svc.UpdateProfile(user)
		assert.NoError(t, err)
	})
}

func TestBooks(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, "key")
	book := &domain.Book{Title: "Title"}

	t.Run("CreateAndList", func(t *testing.T) {
		mockRepo.On("CreateBook", book).Return(nil).Once()
		assert.NoError(t, svc.CreateBook(book))

		mockRepo.On("GetBooks").Return([]domain.Book{*book}, nil).Once()
		res, _ := svc.GetAllBooks()
		assert.Len(t, res, 1)
	})

	t.Run("GetUpdateDelete", func(t *testing.T) {
		mockRepo.On("GetBookByID", uint(1)).Return(book, nil).Once()
		_, err := svc.GetBook(1)
		assert.NoError(t, err)

		mockRepo.On("UpdateBook", book).Return(nil).Once()
		err = svc.UpdateBook(book)
		assert.NoError(t, err)

		mockRepo.On("DeleteBook", uint(1)).Return(nil).Once()
		err = svc.DeleteBook(1)
		assert.NoError(t, err)

		mockRepo.On("GetBooksByAuthor", uint(1)).Return([]domain.Book{}, nil).Once()
		_, err = svc.GetBooksByAuthor(1)
		assert.NoError(t, err)
	})
}

func TestAuthors(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, "key")
	author := &domain.Author{Name: "Author"}

	mockRepo.On("CreateAuthor", author).Return(nil).Once()
	err := svc.CreateAuthor(author)
	assert.NoError(t, err)

	mockRepo.On("GetAuthors").Return([]domain.Author{*author}, nil).Once()
	_, err = svc.GetAllAuthors()
	assert.NoError(t, err)

	mockRepo.On("GetAuthorByID", uint(1)).Return(author, nil).Once()
	_, err = svc.GetAuthor(1)
	assert.NoError(t, err)

	mockRepo.On("UpdateAuthor", author).Return(nil).Once()
	err = svc.UpdateAuthor(author)
	assert.NoError(t, err)

	mockRepo.On("DeleteAuthor", uint(1)).Return(nil).Once()
	err = svc.DeleteAuthor(1)
	assert.NoError(t, err)
}

func TestReviewsAndShelf(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, "key")

	t.Run("Reviews", func(t *testing.T) {
		rev := &domain.Review{Comment: "Good"}
		mockRepo.On("CreateReview", rev).Return(nil).Once()
		err := svc.AddReview(rev)
		assert.NoError(t, err)

		mockRepo.On("GetReviewsByBook", uint(1)).Return([]domain.Review{}, nil).Once()
		_, err = svc.GetReviews(1)
		assert.NoError(t, err)

		mockRepo.On("DeleteReview", uint(1), uint(1)).Return(nil).Once()
		err = svc.DeleteReview(1, 1)
		assert.NoError(t, err)
	})

	t.Run("Shelf", func(t *testing.T) {
		mockRepo.On("AddToShelf", mock.Anything).Return(nil).Once()
		err := svc.SetShelfStatus(1, 1, "reading")
		assert.NoError(t, err)

		mockRepo.On("GetShelf", uint(1)).Return([]domain.Shelf{}, nil).Once()
		_, err = svc.GetShelf(1)
		assert.NoError(t, err)

		mockRepo.On("RemoveFromShelf", uint(1), uint(1)).Return(nil).Once()
		err = svc.RemoveFromShelf(1, 1)
		assert.NoError(t, err)
	})
}
