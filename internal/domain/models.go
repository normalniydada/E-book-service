package domain

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `json:"-"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Author struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"not null" json:"name"`
	Bio   string `gorm:"type:text" json:"bio"`
	Books []Book `gorm:"foreignKey:AuthorID" json:"books,omitempty"`
}

type Book struct {
	ID          uint     `gorm:"primaryKey" json:"id"`
	Title       string   `gorm:"not null" json:"title"`
	Description string   `gorm:"type:text" json:"description"`
	Content     string   `gorm:"type:text" json:"content"`
	AuthorID    uint     `json:"author_id"`
	Author      *Author  `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Reviews     []Review `gorm:"foreignKey:BookID" json:"reviews,omitempty"`
}

type Review struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	BookID  uint   `json:"book_id"`
	UserID  uint   `json:"user_id"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type Shelf struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	BookID    uint      `gorm:"primaryKey" json:"book_id"`
	Book      Book      `gorm:"foreignKey:BookID"`
	Status    string    `json:"status"` // "reading", "completed"
	UpdatedAt time.Time `json:"updated_at"`
}
