package model

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

// Only the hash of the session token is stored so that a leaked table cannot be used to hijack sessions.
type Session struct {
	TokenHash string `gorm:"primarykey"`
	AccountId uint   `gorm:"index;not null"`
	Account   Account
	ExpiresAt time.Time `gorm:"index;not null"`
	CreatedAt time.Time
}

func hashSessionToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// CreateSession stores a new session and returns the raw token to hand to the client.
func CreateSession(ctx context.Context, accountId uint, age time.Duration) (string, error) {
	token := rand.Text()
	session := Session{
		TokenHash: hashSessionToken(token),
		AccountId: accountId,
		ExpiresAt: time.Now().Add(age),
	}
	if err := gorm.G[Session](db).Create(ctx, &session); err != nil {
		return "", err
	}
	return token, nil
}

// FindSession returns the unexpired session for token with its account loaded.
func FindSession(ctx context.Context, token string) (Session, error) {
	return gorm.G[Session](db).
		Preload("Account", nil).
		Where("token_hash = ? AND expires_at > ?", hashSessionToken(token), time.Now()).
		First(ctx)
}

func (session *Session) Extend(ctx context.Context, age time.Duration) error {
	expiresAt := time.Now().Add(age)
	_, err := gorm.G[Session](db).Where("token_hash = ?", session.TokenHash).Update(ctx, "expires_at", expiresAt)
	if err != nil {
		return err
	}
	session.ExpiresAt = expiresAt
	return nil
}

func DeleteSession(ctx context.Context, token string) error {
	_, err := gorm.G[Session](db).Where("token_hash = ?", hashSessionToken(token)).Delete(ctx)
	return err
}

func DeleteExpiredSessions(ctx context.Context) (int, error) {
	return gorm.G[Session](db).Where("expires_at <= ?", time.Now()).Delete(ctx)
}
