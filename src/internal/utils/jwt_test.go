package utils

import (
	"testing"
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/golang-jwt/jwt/v5"
)

func withFixedTime(t *testing.T, initial time.Time) func(time.Time) {
	t.Helper()
	prev := nowFunc
	current := initial
	nowFunc = func() time.Time { return current }
	t.Cleanup(func() { nowFunc = prev })
	return func(next time.Time) { current = next }
}

func TestJWT_ExpiredAccessCanBeRefreshedWhileRefreshValid(t *testing.T) {
	base := time.Date(2026, 11, 26, 12, 0, 0, 0, time.UTC)
	setNow := withFixedTime(t, base)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	accessToken, err := CreateToken(userID, enums.Member)
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}
	refreshToken, err := CreateRefreshToken(userID, enums.Member)
	if err != nil {
		t.Fatalf("CreateRefreshToken failed: %v", err)
	}

	// Advance clock past access token expiry (5h) but within refresh token expiry (7d)
	setNow(base.Add(5*time.Hour + time.Minute))

	if _, err := VerifyToken(accessToken); err == nil {
		t.Fatalf("expected access token to be expired after 5h1m")
	}

	newAccess, newRefresh, err := RefreshAccessToken(refreshToken)
	if err != nil {
		t.Fatalf("expected RefreshAccessToken to succeed while refresh token is valid, got: %v", err)
	}

	accessClaims, err := VerifyToken(newAccess)
	if err != nil {
		t.Fatalf("expected new access token to be valid, got: %v", err)
	}
	if accessClaims["type"] != "access" || accessClaims["sub"] != userID {
		t.Fatalf("unexpected access claims: %+v", accessClaims)
	}

	refreshClaims, err := VerifyToken(newRefresh)
	if err != nil {
		t.Fatalf("expected new refresh token to be valid, got: %v", err)
	}
	if refreshClaims["type"] != "refresh" || refreshClaims["sub"] != userID {
		t.Fatalf("unexpected refresh claims: %+v", refreshClaims)
	}
}

func TestJWT_RefreshAccessTokenRejectsExpiredRefreshToken(t *testing.T) {
	base := time.Date(2026, 11, 26, 12, 0, 0, 0, time.UTC)
	setNow := withFixedTime(t, base)

	refreshToken, err := CreateRefreshToken("550e8400-e29b-41d4-a716-446655440000", enums.Member)
	if err != nil {
		t.Fatalf("CreateRefreshToken failed: %v", err)
	}

	// Advance clock past refresh token expiry (7d)
	setNow(base.Add(7*24*time.Hour + time.Minute))

	if _, _, err := RefreshAccessToken(refreshToken); err == nil {
		t.Fatalf("expected RefreshAccessToken to reject expired refresh token")
	}
}

func TestJWT_RefreshAccessTokenRejectsAccessToken(t *testing.T) {
	base := time.Date(2026, 11, 26, 12, 0, 0, 0, time.UTC)
	withFixedTime(t, base)

	accessToken, err := CreateToken("550e8400-e29b-41d4-a716-446655440000", enums.Member)
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}

	if _, _, err := RefreshAccessToken(accessToken); err == nil {
		t.Fatalf("expected RefreshAccessToken to reject an access token")
	}
}

func TestJWT_VerifyTokenRejectsMissingExpirationAndWrongAlg(t *testing.T) {
	base := time.Date(2026, 11, 26, 12, 0, 0, 0, time.UTC)
	withFixedTime(t, base)

	// Missing exp claim
	noExp := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user-1",
		"role": string(enums.Member),
		"type": "access",
	})
	noExpStr, err := noExp.SignedString(secretKey)
	if err != nil {
		t.Fatalf("sign noExp failed: %v", err)
	}
	if _, err := VerifyToken(noExpStr); err == nil {
		t.Fatalf("expected VerifyToken to reject token without exp claim")
	}

	// Wrong signing method (HS512 instead of HS256)
	wrongAlg := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":  "user-1",
		"role": string(enums.Member),
		"type": "access",
		"exp":  base.Add(time.Hour).Unix(),
	})
	wrongAlgStr, err := wrongAlg.SignedString(secretKey)
	if err != nil {
		t.Fatalf("sign wrongAlg failed: %v", err)
	}
	if _, err := VerifyToken(wrongAlgStr); err == nil {
		t.Fatalf("expected VerifyToken to reject HS512 token")
	}
}
