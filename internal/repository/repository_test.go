package repository

import (
	"E-book-service/internal/domain"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepoTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo Repository
}

func (s *RepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(s.T(), err)

	s.mock = mock
	s.repo = NewRepository(gormDB)
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}

// --- USERS ---

func (s *RepoTestSuite) TestUsers() {
	user := &domain.User{Email: "test@test.com", Name: "Name"}
	user.ID = 1

	// CreateUser
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()
	assert.NoError(s.T(), s.repo.CreateUser(user))

	// GetUserByEmail
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
		WithArgs("test@test.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "test@test.com"))
	_, err := s.repo.GetUserByEmail("test@test.com")
	assert.NoError(s.T(), err)

	// GetUserByID
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = s.repo.GetUserByID(1)
	assert.NoError(s.T(), err)

	// UpdateUser
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err = s.repo.UpdateUser(user)
	assert.NoError(s.T(), err)
}

// --- BOOKS ---

func (s *RepoTestSuite) TestBooks() {
	book := &domain.Book{Title: "Title", AuthorID: 1}
	book.ID = 1

	// CreateBook
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "books"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()
	err := s.repo.CreateBook(book)
	assert.NoError(s.T(), err)

	// GetBooks
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "author_id"}).AddRow(1, 1))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "authors" WHERE "authors"."id" = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = s.repo.GetBooks()
	assert.NoError(s.T(), err)

	// GetBookByID
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE "books"."id" = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "author_id"}).AddRow(1, 1))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "authors" WHERE "authors"."id" = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = s.repo.GetBookByID(1)
	assert.NoError(s.T(), err)

	// UpdateBook
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET`)).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err = s.repo.UpdateBook(book)
	assert.NoError(s.T(), err)

	// DeleteBook
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "books" WHERE "books"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err = s.repo.DeleteBook(1)
	assert.NoError(s.T(), err)

	// GetBooksByAuthor
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE author_id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = s.repo.GetBooksByAuthor(1)
	assert.NoError(s.T(), err)
}

// --- AUTHORS ---

func (s *RepoTestSuite) TestAuthors() {
	author := &domain.Author{Name: "Name"}
	author.ID = 1

	// CreateAuthor
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "authors"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()
	err := s.repo.CreateAuthor(author)
	assert.NoError(s.T(), err)

	// GetAuthors
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "authors"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = s.repo.GetAuthors()
	assert.NoError(s.T(), err)

	// GetAuthorByID
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "authors" WHERE "authors"."id" = $1`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = s.repo.GetAuthorByID(1)
	assert.NoError(s.T(), err)

	// UpdateAuthor
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "authors" SET`)).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err = s.repo.UpdateAuthor(author)
	assert.NoError(s.T(), err)

	// DeleteAuthor
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "authors" WHERE "authors"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err = s.repo.DeleteAuthor(1)
	assert.NoError(s.T(), err)
}

// --- REVIEWS ---

func (s *RepoTestSuite) TestReviews() {
	review := &domain.Review{Comment: "C", BookID: 1, UserID: 1}

	// CreateReview
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "reviews"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()
	err := s.repo.CreateReview(review)
	assert.NoError(s.T(), err)

	// GetReviewsByBook
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reviews" WHERE book_id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = s.repo.GetReviewsByBook(1)
	assert.NoError(s.T(), err)

	// DeleteReview
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "reviews" WHERE id = $1 AND user_id = $2`)).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err = s.repo.DeleteReview(1, 1)
	assert.NoError(s.T(), err)
}

// --- SHELF ---

func (s *RepoTestSuite) TestShelf() {
	shelf := &domain.Shelf{UserID: 1, BookID: 1, Status: "reading", UpdatedAt: time.Now()}

	// AddToShelf (Save)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "shelves" SET`)).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.repo.AddToShelf(shelf)
	assert.NoError(s.T(), err)

	// GetShelf
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "shelves" WHERE user_id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "book_id"}).AddRow(1, 1))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE "books"."id" = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "author_id"}).AddRow(1, 1))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "authors" WHERE "authors"."id" = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = s.repo.GetShelf(1)
	assert.NoError(s.T(), err)

	// RemoveFromShelf
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "shelves" WHERE user_id = $1 AND book_id = $2`)).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err = s.repo.RemoveFromShelf(1, 1)
	assert.NoError(s.T(), err)
}
