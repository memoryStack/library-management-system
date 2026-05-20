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

// JWT validates access tokens for the primary Auth0 application (Conf).
var JWT *validator.Validator

// JWTPasswordless validates access tokens for the passwordless Auth0 application, when configured.
var JWTPasswordless *validator.Validator

func initJWT(cfg *Config) error {
	v, err := newJWTValidator(cfg)
	if err != nil {
		return err
	}
	JWT = v
	return nil
}

func initJWTPasswordless(cfg *Config) error {
	if cfg == nil {
		return nil
	}
	if Conf != nil && cfg.Domain == Conf.Domain && cfg.Audience == Conf.Audience {
		JWTPasswordless = JWT
		return nil
	}
	v, err := newJWTValidator(cfg)
	if err != nil {
		return err
	}
	JWTPasswordless = v
	return nil
}

func newJWTValidator(cfg *Config) (*validator.Validator, error) {
	issuer := "https://" + cfg.Domain + "/"
	issuerURL, err := url.Parse(issuer)
	if err != nil {
		return nil, fmt.Errorf("issuer url: %w", err)
	}
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	return validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuer,
		[]string{cfg.Audience},
		validator.WithAllowedClockSkew(30*time.Second),
	)
}

func validatorFor(cfg *Config) *validator.Validator {
	if cfg == nil {
		return JWT
	}
	if PasswordlessConf != nil && cfg.ClientID == PasswordlessConf.ClientID {
		return JWTPasswordless
	}
	return JWT
}

// AuthConfigs returns primary and passwordless configs (non-nil entries only).
func AuthConfigs() []*Config {
	var out []*Config
	if Conf != nil {
		out = append(out, Conf)
	}
	if PasswordlessConf != nil {
		out = append(out, PasswordlessConf)
	}
	return out
}

// ValidateAccessToken validates the token with the primary Auth0 validator.
func ValidateAccessToken(ctx context.Context, raw string) (*validator.ValidatedClaims, error) {
	return validateWith(JWT, raw, ctx)
}

// ValidateAccessTokenForConfig validates using the validator for the given Auth0 app config.
func ValidateAccessTokenForConfig(ctx context.Context, cfg *Config, raw string) (*validator.ValidatedClaims, error) {
	return validateWith(validatorFor(cfg), raw, ctx)
}

// ValidateAccessTokenAny tries each configured Auth0 app until one validates the token.
func ValidateAccessTokenAny(ctx context.Context, raw string) (*validator.ValidatedClaims, *Config, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil, fmt.Errorf("empty token")
	}
	var lastErr error
	for _, cfg := range AuthConfigs() {
		claims, err := ValidateAccessTokenForConfig(ctx, cfg, raw)
		if err == nil {
			return claims, cfg, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return nil, nil, lastErr
	}
	return nil, nil, fmt.Errorf("jwt validator not initialized")
}

func validateWith(v *validator.Validator, raw string, ctx context.Context) (*validator.ValidatedClaims, error) {
	if v == nil {
		return nil, fmt.Errorf("jwt validator not initialized")
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty token")
	}
	claims, err := v.ValidateToken(ctx, raw)
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
