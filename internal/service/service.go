package service

import (
	"E-book-service/internal/domain"
	"E-book-service/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type ServiceInterface interface {
	Register(email, pass, name string) error
	Login(email, pass string) (string, error)
	GetProfile(id uint) (*domain.User, error)
	UpdateProfile(u *domain.User) error
	CreateBook(b *domain.Book) error
	GetAllBooks() ([]domain.Book, error)
	GetBook(id uint) (*domain.Book, error)
	UpdateBook(b *domain.Book) error
	DeleteBook(id uint) error
	GetBooksByAuthor(aID uint) ([]domain.Book, error)
	CreateAuthor(a *domain.Author) error
	GetAllAuthors() ([]domain.Author, error)
	GetAuthor(id uint) (*domain.Author, error)
	UpdateAuthor(a *domain.Author) error
	DeleteAuthor(id uint) error
	AddReview(re *domain.Review) error
	GetReviews(bID uint) ([]domain.Review, error)
	DeleteReview(id, uID uint) error
	SetShelfStatus(uID, bID uint, status string) error
	GetShelf(uID uint) ([]domain.Shelf, error)
	RemoveFromShelf(uID, bID uint) error
}

type service struct {
	repo   repository.Repository // Используем интерфейс!
	jwtKey string
}

func NewService(r repository.Repository, key string) ServiceInterface {
	return &service{repo: r, jwtKey: key}
}

func (s *service) Register(email, pass, name string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return s.repo.CreateUser(&domain.User{Email: email, Password: string(hash), Name: name})
}

func (s *service) Login(email, pass string) (string, error) {
	u, err := s.repo.GetUserByEmail(email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass)) != nil {
		return "", errors.New("invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  u.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(s.jwtKey))
}

func (s *service) GetProfile(id uint) (*domain.User, error) { return s.repo.GetUserByID(id) }
func (s *service) UpdateProfile(u *domain.User) error       { return s.repo.UpdateUser(u) }

// BOOKS
func (s *service) CreateBook(b *domain.Book) error       { return s.repo.CreateBook(b) }
func (s *service) GetAllBooks() ([]domain.Book, error)   { return s.repo.GetBooks() }
func (s *service) GetBook(id uint) (*domain.Book, error) { return s.repo.GetBookByID(id) }
func (s *service) UpdateBook(b *domain.Book) error       { return s.repo.UpdateBook(b) }
func (s *service) DeleteBook(id uint) error              { return s.repo.DeleteBook(id) }
func (s *service) GetBooksByAuthor(aID uint) ([]domain.Book, error) {
	return s.repo.GetBooksByAuthor(aID)
}

// AUTHORS
func (s *service) CreateAuthor(a *domain.Author) error       { return s.repo.CreateAuthor(a) }
func (s *service) GetAllAuthors() ([]domain.Author, error)   { return s.repo.GetAuthors() }
func (s *service) GetAuthor(id uint) (*domain.Author, error) { return s.repo.GetAuthorByID(id) }
func (s *service) UpdateAuthor(a *domain.Author) error       { return s.repo.UpdateAuthor(a) }
func (s *service) DeleteAuthor(id uint) error                { return s.repo.DeleteAuthor(id) }

// REVIEWS
func (s *service) AddReview(re *domain.Review) error            { return s.repo.CreateReview(re) }
func (s *service) GetReviews(bID uint) ([]domain.Review, error) { return s.repo.GetReviewsByBook(bID) }
func (s *service) DeleteReview(id, uID uint) error              { return s.repo.DeleteReview(id, uID) }

// SHELF
func (s *service) SetShelfStatus(uID, bID uint, status string) error {
	return s.repo.AddToShelf(&domain.Shelf{UserID: uID, BookID: bID, Status: status, UpdatedAt: time.Now()})
}
func (s *service) GetShelf(uID uint) ([]domain.Shelf, error) { return s.repo.GetShelf(uID) }
func (s *service) RemoveFromShelf(uID, bID uint) error       { return s.repo.RemoveFromShelf(uID, bID) }
