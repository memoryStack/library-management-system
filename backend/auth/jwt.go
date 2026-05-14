package auth

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	gojwt "gopkg.in/go-jose/go-jose.v2/jwt"
)

// JWT is used by middleware to validate access tokens (signature, iss, aud, exp).
var JWT *validator.Validator

func initJWT(cfg *Config) error {
	issuer := "https://" + cfg.Domain + "/"
	issuerURL, err := url.Parse(issuer)
	if err != nil {
		return fmt.Errorf("issuer url: %w", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	v, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuer,
		[]string{cfg.Audience},
		validator.WithAllowedClockSkew(30*time.Second),
	)
	if err != nil {
		return fmt.Errorf("jwt validator: %w", err)
	}
	JWT = v
	return nil
}

// ValidateAccessToken validates the bearer access JWT (API audience).
func ValidateAccessToken(ctx context.Context, raw string) (*validator.ValidatedClaims, error) {
	if JWT == nil {
		return nil, fmt.Errorf("jwt validator not initialized")
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty token")
	}
	claims, err := JWT.ValidateToken(ctx, raw)
	if err != nil {
		return nil, err
	}
	validated, ok := claims.(*validator.ValidatedClaims)
	if !ok || validated == nil {
		return nil, fmt.Errorf("unexpected claims type")
	}
	return validated, nil
}

// IDTokenClaims returns all id_token claims as interface{} (map[string]interface{}).
// This only decodes claims; add signature/iss/aud validation if used for auth decisions.
func IDTokenClaims(raw string) (map[string]interface{}, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]interface{}{}, fmt.Errorf("empty id_token")
	}
	tok, err := gojwt.ParseSigned(raw)
	if err != nil {
		return nil, fmt.Errorf("parse id_token: %w", err)
	}
	claims := map[string]interface{}{}
	if err := tok.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return nil, fmt.Errorf("decode id_token claims: %w", err)
	}
	return claims, nil
}
