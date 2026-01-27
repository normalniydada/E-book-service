package repository

import (
	"E-book-service/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	// Users
	CreateUser(u *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	UpdateUser(u *domain.User) error

	// Books
	CreateBook(b *domain.Book) error
	GetBooks() ([]domain.Book, error)
	GetBookByID(id uint) (*domain.Book, error)
	UpdateBook(b *domain.Book) error
	DeleteBook(id uint) error
	GetBooksByAuthor(aID uint) ([]domain.Book, error)

	// Authors
	CreateAuthor(a *domain.Author) error
	GetAuthors() ([]domain.Author, error)
	GetAuthorByID(id uint) (*domain.Author, error)
	UpdateAuthor(a *domain.Author) error
	DeleteAuthor(id uint) error

	// Reviews
	CreateReview(re *domain.Review) error
	GetReviewsByBook(bookID uint) ([]domain.Review, error)
	DeleteReview(id, uID uint) error

	// Shelf
	AddToShelf(s *domain.Shelf) error
	GetShelf(uID uint) ([]domain.Shelf, error)
	RemoveFromShelf(uID, bID uint) error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateUser(u *domain.User) error { return r.db.Create(u).Error }
func (r *postgresRepository) GetUserByEmail(email string) (*domain.User, error) {
	var u domain.User
	return &u, r.db.Where("email = ?", email).First(&u).Error
}
func (r *postgresRepository) GetUserByID(id uint) (*domain.User, error) {
	var u domain.User
	return &u, r.db.First(&u, id).Error
}
func (r *postgresRepository) UpdateUser(u *domain.User) error { return r.db.Save(u).Error }

func (r *postgresRepository) CreateBook(b *domain.Book) error { return r.db.Create(b).Error }
func (r *postgresRepository) GetBooks() ([]domain.Book, error) {
	var b []domain.Book
	return b, r.db.Preload("Author").Find(&b).Error
}
func (r *postgresRepository) GetBookByID(id uint) (*domain.Book, error) {
	var b domain.Book
	return &b, r.db.Preload("Author").First(&b, id).Error
}
func (r *postgresRepository) UpdateBook(b *domain.Book) error { return r.db.Save(b).Error }
func (r *postgresRepository) DeleteBook(id uint) error        { return r.db.Delete(&domain.Book{}, id).Error }
func (r *postgresRepository) GetBooksByAuthor(aID uint) ([]domain.Book, error) {
	var b []domain.Book
	return b, r.db.Where("author_id = ?", aID).Find(&b).Error
}

func (r *postgresRepository) CreateAuthor(a *domain.Author) error { return r.db.Create(a).Error }
func (r *postgresRepository) GetAuthors() ([]domain.Author, error) {
	var a []domain.Author
	return a, r.db.Find(&a).Error
}
func (r *postgresRepository) GetAuthorByID(id uint) (*domain.Author, error) {
	var a domain.Author
	return &a, r.db.First(&a, id).Error
}
func (r *postgresRepository) UpdateAuthor(a *domain.Author) error { return r.db.Save(a).Error }
func (r *postgresRepository) DeleteAuthor(id uint) error {
	return r.db.Delete(&domain.Author{}, id).Error
}

func (r *postgresRepository) CreateReview(re *domain.Review) error { return r.db.Create(re).Error }
func (r *postgresRepository) GetReviewsByBook(bookID uint) ([]domain.Review, error) {
	var re []domain.Review
	return re, r.db.Where("book_id = ?", bookID).Find(&re).Error
}
func (r *postgresRepository) DeleteReview(id, uID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, uID).Delete(&domain.Review{}).Error
}

func (r *postgresRepository) AddToShelf(s *domain.Shelf) error { return r.db.Save(s).Error }
func (r *postgresRepository) GetShelf(uID uint) ([]domain.Shelf, error) {
	var s []domain.Shelf
	return s, r.db.Preload("Book.Author").Where("user_id = ?", uID).Find(&s).Error
}
func (r *postgresRepository) RemoveFromShelf(uID, bID uint) error {
	return r.db.Where("user_id = ? AND book_id = ?", uID, bID).Delete(&domain.Shelf{}).Error
}
