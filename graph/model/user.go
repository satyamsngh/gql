package model

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"hash"
	"os"
	"strconv"
	"time"
)

var (
	PwPepper              = os.Getenv("PASSWORD_PEPPER")
	ErrInvalidCredentials = errors.New("Invalid username or password")
	RememberTokenLength   = 32
)

// Conn is our main struct, including the database instance for working with data.
type UserService struct {
	db   *gorm.DB
	hmac hash.Hash
}

func NewUserService(db *gorm.DB) *UserService {
	hmacSecret := os.Getenv("HMAC_SECRET_KEY")
	return &UserService{
		db:   db,
		hmac: hmac.New(sha256.New, []byte(hmacSecret)),
	}
}

type User struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Remember     string
}

// NewService is the constructor for the Conn struct.
func NewService(db *gorm.DB) (*UserService, error) {

	// We check if the database instance is nil, which would indicate an issue.
	if db == nil {
		return nil, errors.New("please provide a valid connection")
	}

	// We initialize our service with the passed database instance.
	s := &UserService{db: db}
	return s, nil
}

// CreateUser is a method that creates a new user record in the database.
func (s *UserService) CreateUser(input NewUser) (*User, error) {

	// We hash the user's password for storage in the database.
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generating password hash: %w", err)
	}

	// We prepare the User record.
	u := &User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashedPass),
	}

	// We attempt to create the new User record in the database.
	err = s.db.Create(&u).Error
	if err != nil {
		return nil, err
	}

	// Successfully created the record, return the user.
	return u, nil
}
func (u *UserService) GetBy(field string, val interface{}) (*User, error) {
	var user User
	result := u.db.Where(fmt.Sprintf("%s = ?", field), val).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, result.Error
	}
	return &user, nil
}

// Authenticate is a method that checks a user's provided email and password against the database.
func (s *UserService) Authenticate(ctx context.Context, email, password string) (jwt.RegisteredClaims,
	error) {

	// We attempt to find the User record where the email
	// matches the provided email.
	var u User
	tx := s.db.Where("email = ?", email).First(&u)
	if tx.Error != nil {
		return jwt.RegisteredClaims{}, tx.Error
	}

	// We check if the provided password matches the hashed password in the database.
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}

	// Successful authentication! Generate JWT claims.
	c := jwt.RegisteredClaims{
		Issuer:    "graphql",
		Subject:   strconv.FormatUint(uint64(u.ID), 10),
		Audience:  jwt.ClaimStrings{"students"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	// And return those claims.
	return c, nil
}
