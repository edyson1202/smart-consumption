package models

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

var IncorrectCredentialsError = errors.New("Incorrect username or password")

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type NewUser struct {
	Email    string
	Password string
}

type DataPoint struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Consumption float64 `json:"consumption"`
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(newUser NewUser) (*User, error) {
	newUser.Email = strings.ToLower(newUser.Email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{
		Email:        newUser.Email,
		PasswordHash: passwordHash,
	}
	row := us.DB.QueryRow(`
INSERT INTO USERS(email, password_hash)
VALUES ($1, $2) RETURNING id`, newUser.Email, passwordHash)

	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}

	return &user, nil
}

func (us *UserService) Authenticate(user NewUser) (*User, error) {
	// TODO this is done both in the Create and Authenticate method, think of a way to do this validation in one place avoiding duplication
	user.Email = strings.ToLower(user.Email)

	row := us.DB.QueryRow(`
SELECT id, email, password_hash FROM USERS WHERE email = $1`, user.Email)

	u := User{}
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("Authenticate user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(user.Password))
	if err != nil {
		return nil, IncorrectCredentialsError
	}

	return &u, nil
}

func (us *UserService) CreateDatapoint(userId int, dataPoint DataPoint) error {

	query := `
INSERT INTO DATAPOINT(user_id, measured_at, temperature, humidity, consumption)
VALUES ($1, NOW(), $2, $3, $4) RETURNING id`

	_, err := us.DB.Exec(query, userId, dataPoint.Temperature, dataPoint.Humidity, dataPoint.Consumption)

	return err
}
