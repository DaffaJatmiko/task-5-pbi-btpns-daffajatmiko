package models

import (
	"errors"
	"html"
	"strings"
	"task-5-pbi-btpns-daffajatmiko/app"
	"task-5-pbi-btpns-daffajatmiko/helpers/hash"
	"time"

	"github.com/badoux/checkmail"
	"github.com/google/uuid"
)

// User represents the user model.
type User struct {
	ID        string    `gorm:"primary_key; unique" json:"id"`
	Username  string    `gorm:"size:255;not null;" json:"username"`
	Email     string    `gorm:"size:255;not null; unique" json:"email"`
	Password  string    `gorm:"size:255;not null;" json:"password"`
	Photos    Photo     `gorm:"constraint:OnUpdate:CASCADE, OnDelete:SET NULL;" json:"photos"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Photo represents the photo model.
type Photo struct {
	ID       int       `gorm:"primary_key;auto_increment" json:"id"`
	Title    string    `gorm:"size:255;not null" json:"title"`
	Caption  string    `gorm:"size:255;not null" json:"caption"`
	PhotoUrl string    `gorm:"size:255;not null;" json:"photo_url"`
	UserID   string    `gorm:"not null" json:"user_id"`
	Owner    app.Owner `gorm:"owner"`
}

// USER METHODS

// Init initializes user data.
func (u *User) Init() {
	u.ID = uuid.New().String()                                    // Generate a new UUID
	u.Username = html.EscapeString(strings.TrimSpace(u.Username)) // Escape string
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
}

// HashPassword changes the password to a hashed password.
func (u *User) HashPassword() error {
	hashedPassword, err := hash.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks the provided password.
func (u *User) CheckPassword(providedPassword string) error {
	err := hash.CheckPasswordHash(u.Password, providedPassword)
	if err != nil {
		return err
	}
	return nil
}

// Validate validates user data based on the given action.
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Email == "" {
			return errors.New("email is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("email is invalid")
		}
		return nil

	case "register":
		if u.ID == "" {
			return errors.New("id is required")
		} else if u.Email == "" {
			return errors.New("email is required")
		} else if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("email is invalid")
		} else if u.Username == "" {
			return errors.New("username is required")
		} else if u.Password == "" {
			return errors.New("password is required")
		} else if len(u.Password) < 8 {
			return errors.New("password must be at least 8 characters")
		}
		return nil

	case "update":
		if u.ID == "" {
			return errors.New("id is required")
		} else if u.Email == "" {
			return errors.New("email is required")
		} else if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		} else if u.Username == "" {
			return errors.New("username is required")
		} else if u.Password == "" {
			return errors.New("password is required")
		} else if len(u.Password) < 8 {
			return errors.New("password must be at least 8 characters")
		}
		return nil

	default:
		return nil
	}
}

// PHOTO METHODS

// Init initializes Photo data.
func (p *Photo) Init() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title)) // Escape string
	p.Caption = html.EscapeString(strings.TrimSpace(p.Caption))
	p.PhotoUrl = html.EscapeString(strings.TrimSpace(p.PhotoUrl))
}

// Validate validates Photo data based on the given action.
func (p *Photo) Validate(action string) error {
	switch strings.ToLower(action) {
	case "upload":
		if p.Title == "" {
			return errors.New("title is required")
		} else if p.Caption == "" {
			return errors.New("caption is required")
		} else if p.UserID == "" {
			return errors.New("userID is required")
		}
		return nil

	case "change":
		if p.Title == "" {
			return errors.New("title is required")
		} else if p.Caption == "" {
			return errors.New("caption is required")
		} else if p.PhotoUrl == "" {
			return errors.New("photoUrl is required")
		}
		return nil

	default:
		return nil
	}
}
