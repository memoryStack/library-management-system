package auth

import (
	"context"
	"fmt"
)

// Init loads Auth0 configuration and sets up the JWT validator (JWKS).
func Init(ctx context.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	Conf = cfg
	if err := initJWT(cfg); err != nil {
		return fmt.Errorf("jwt: %w", err)
	}
	_ = ctx
	return nil
}
