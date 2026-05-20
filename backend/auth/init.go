package auth

import (
	"context"
	"fmt"
)

// Init loads primary and optional passwordless Auth0 configuration and sets up the JWT validator (JWKS).
// JWT validation uses the primary Conf (tenant/audience); if passwordless uses a different tenant or API,
// extend initJWT to validate multiple issuers/audiences as needed.
func Init(ctx context.Context) error {
	cfg, err := loadAuth0Config("AUTH0_")
	if err != nil {
		return err
	}
	Conf = cfg

	pl, err := loadPasswordlessConfigOptional()
	if err != nil {
		return err
	}
	PasswordlessConf = pl

	if err := initJWT(cfg); err != nil {
		return fmt.Errorf("jwt: %w", err)
	}
	if err := initJWTPasswordless(pl); err != nil {
		return fmt.Errorf("jwt passwordless: %w", err)
	}
	_ = ctx
	return nil
}
