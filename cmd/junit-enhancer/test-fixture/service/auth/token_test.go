package auth

import (
	"testing"
	"time"
)

func TestTokenManager_GenerateToken(t *testing.T) {
	tm := NewTokenManager()
	
	token, err := tm.GenerateToken(123, time.Hour)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if token.UserID != 123 {
		t.Errorf("Expected UserID 123, got %d", token.UserID)
	}
	
	if token.Value == "" {
		t.Error("Expected non-empty token value")
	}
	
	if time.Now().After(token.ExpiresAt) {
		t.Error("Token should not be expired immediately")
	}
}

func TestTokenManager_ValidateToken(t *testing.T) {
	tm := NewTokenManager()
	
	t.Run("valid_token", func(t *testing.T) {
		generated, err := tm.GenerateToken(456, time.Hour)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}
		
		validated, err := tm.ValidateToken(generated.Value)
		if err != nil {
			t.Fatalf("Token validation failed: %v", err)
		}
		
		if validated.UserID != 456 {
			t.Errorf("Expected UserID 456, got %d", validated.UserID)
		}
	})
	
	t.Run("invalid_token", func(t *testing.T) {
		_, err := tm.ValidateToken("invalid-token")
		if err == nil {
			t.Error("Expected error for invalid token")
		}
	})
	
	t.Run("expired_token", func(t *testing.T) {
		generated, err := tm.GenerateToken(789, -time.Hour) // Already expired
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}
		
		_, err = tm.ValidateToken(generated.Value)
		if err == nil {
			t.Error("Expected error for expired token")
		}
	})
}

func TestTokenManager_RevokeToken(t *testing.T) {
	tm := NewTokenManager()
	
	t.Run("existing_token", func(t *testing.T) {
		token, err := tm.GenerateToken(111, time.Hour)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}
		
		revoked := tm.RevokeToken(token.Value)
		if !revoked {
			t.Error("Expected token to be revoked")
		}
		
		// Verify token is no longer valid
		_, err = tm.ValidateToken(token.Value)
		if err == nil {
			t.Error("Expected error for revoked token")
		}
	})
	
	t.Run("non_existent_token", func(t *testing.T) {
		revoked := tm.RevokeToken("non-existent")
		if revoked {
			t.Error("Expected revocation to fail for non-existent token")
		}
	})
}

func FuzzGenerateToken(f *testing.F) {
	f.Add(1, int64(time.Hour))
	f.Add(999999, int64(time.Minute))
	f.Add(-1, int64(time.Second))
	
	f.Fuzz(func(t *testing.T, userID int, durationNanos int64) {
		if durationNanos <= 0 {
			return // Skip negative or zero durations
		}
		
		tm := NewTokenManager()
		duration := time.Duration(durationNanos)
		
		token, err := tm.GenerateToken(userID, duration)
		if err != nil {
			t.Errorf("Unexpected error generating token: %v", err)
			return
		}
		
		if token.UserID != userID {
			t.Errorf("Expected UserID %d, got %d", userID, token.UserID)
		}
		
		if token.Value == "" {
			t.Error("Token value should not be empty")
		}
		
		// Validate the token immediately
		validated, err := tm.ValidateToken(token.Value)
		if err != nil {
			t.Errorf("Token should be valid immediately after generation: %v", err)
		}
		
		if validated.UserID != userID {
			t.Errorf("Validated token UserID mismatch: expected %d, got %d", userID, validated.UserID)
		}
	})
}
