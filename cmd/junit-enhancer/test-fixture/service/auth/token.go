package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

type Token struct {
	Value     string
	UserID    int
	ExpiresAt time.Time
}

type TokenManager struct {
	tokens map[string]*Token
}

func NewTokenManager() *TokenManager {
	return &TokenManager{
		tokens: make(map[string]*Token),
	}
}

func (tm *TokenManager) GenerateToken(userID int, duration time.Duration) (*Token, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	
	token := &Token{
		Value:     hex.EncodeToString(bytes),
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
	}
	
	tm.tokens[token.Value] = token
	return token, nil
}

func (tm *TokenManager) ValidateToken(tokenValue string) (*Token, error) {
	token, exists := tm.tokens[tokenValue]
	if !exists {
		return nil, errors.New("token not found")
	}
	
	if time.Now().After(token.ExpiresAt) {
		delete(tm.tokens, tokenValue)
		return nil, errors.New("token expired")
	}
	
	return token, nil
}

func (tm *TokenManager) RevokeToken(tokenValue string) bool {
	if _, exists := tm.tokens[tokenValue]; exists {
		delete(tm.tokens, tokenValue)
		return true
	}
	return false
}
