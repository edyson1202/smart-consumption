package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"lenslocked/pkg/rand"
)

const (
	minBytesPerToken   = 32
	sessionTokenLength = 32
)

type Session struct {
	Id        int
	UserId    int
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userId int) (*Session, string, error) {
	token := rand.String(sessionTokenLength)

	session := &Session{
		UserId:    userId,
		TokenHash: ss.hash(token),
	}

	row := ss.DB.QueryRow(`
INSERT INTO session (user_id, token_hash)
VALUES($1, $2) ON CONFLICT (user_id) DO
UPDATE
SET token_hash = $2
RETURNING id`, session.UserId, session.TokenHash)

	err := row.Scan(&session.Id)
	if err != nil {
		return nil, "", fmt.Errorf("Create: %w", err)
	}
	return session, token, nil
}

func (ss *SessionService) User(token string) (*User, error) {

	tokenHash := ss.hash(token)

	var user User
	row := ss.DB.QueryRow(`SELECT users.id, users.email, users.password_hash FROM users JOIN session
ON users.id = session.user_id WHERE session.token_hash = $1`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	_, err := ss.DB.Exec(`DELETE FROM session WHERE token_hash = $1`, ss.hash(token))
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
