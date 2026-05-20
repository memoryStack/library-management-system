package controllers

import (
	"crypto/rand"
	"crypto/subtle"
	"bytes"
	"encoding/json"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"library-management-system/backend/auth"

	"github.com/gofiber/fiber/v2"

	"library-management-system/backend/initializers"
	"library-management-system/backend/models"

	"gorm.io/gorm"
)

func cookieSameSite(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "strict":
		return fiber.CookieSameSiteStrictMode
	case "none":
		return fiber.CookieSameSiteNoneMode
	default:
		return fiber.CookieSameSiteLaxMode
	}
}

func setStateCookie(c *fiber.Ctx, cfg *auth.Config, state string) {
	c.Cookie(&fiber.Cookie{
		Name:     cfg.StateCookieName,
		Value:    state,
		Path:     "/api/auth",
		HTTPOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: cookieSameSite(cfg.CookieSameSite),
		MaxAge:   600,
	})
}

func clearStateCookie(c *fiber.Ctx, cfg *auth.Config) {
	c.Cookie(&fiber.Cookie{
		Name:     cfg.StateCookieName,
		Value:    "",
		Path:     "/api/auth",
		HTTPOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: cookieSameSite(cfg.CookieSameSite),
		MaxAge:   -1,
	})
}

func setAuthCookies(c *fiber.Ctx, cfg *auth.Config, tr *auth.TokenResponse) {
	auth.SetAuthCookies(c, cfg, tr)
}

func clearAuthCookies(c *fiber.Ctx, cfg *auth.Config) {
	ss := cookieSameSite(cfg.CookieSameSite)
	clear := func(name, path string) {
		c.Cookie(&fiber.Cookie{
			Name:     name,
			Value:    "",
			Path:     path,
			HTTPOnly: true,
			Secure:   cfg.CookieSecure,
			SameSite: ss,
			MaxAge:   -1,
		})
	}
	clear(cfg.AccessCookieName, "/")
	clear(cfg.RefreshCookieName, cfg.RefreshCookiePath)
	clear(cfg.RefreshCookieName, "/")
}

// AuthLogin starts the Auth0 Universal Login flow (passwordless is configured in Auth0).
func AuthLogin(c *fiber.Ctx) error {
	fmt.Println("AuthLogin started")
	cfg := auth.GetAuthConfigs(c)
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create login state"})
	}
	state := hex.EncodeToString(b)
	setStateCookie(c, cfg, state)
	fmt.Println("AuthLogin started")
	return c.Redirect(auth.AuthorizeURL(cfg, state), fiber.StatusFound)
}

func AuthLoginSelfManaged(c *fiber.Ctx) error {

	cfg := auth.PasswordlessConf
	
	url := "https://"+auth.PasswordlessConf.Domain+"/passwordless/start"
	
	requestPayload := map[string]string{
		"client_id": cfg.ClientID,
		"client_secret": cfg.ClientSecret,
		"connection": c.Query("medium"),
		"send": "code",
	}

	var requestBody map[string]string
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if c.Query("medium") == "email" {
		requestPayload["email"] = requestBody["email"]
	} else {
		// phone number with code
		requestPayload["phone_number"] = requestBody["phone_number"]
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Passwordless login started",
	})
}

func AuthConfirmOTP(c *fiber.Ctx) error {
	if c.Locals("clientKnown") != true {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not determine client type; pass client=web or client=native (query or X-Client-Type header)",
		})
	}

	cfg := auth.PasswordlessConf
	url := "https://" + cfg.Domain + "/oauth/token"

	var requestBody map[string]string
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	requestPayload := map[string]string{
		"grant_type":    "http://auth0.com/oauth/grant-type/passwordless/otp",
		"client_id":     cfg.ClientID,
		"client_secret": cfg.ClientSecret,
		"audience":      cfg.Audience,
		"username":      requestBody["username"],
		"otp":           requestBody["otp"],
		"realm":         c.Query("medium"),
		"scope":         "openid profile email offline_access",
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error": string(body),
		})
	}

	var tr auth.TokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println("access token: ", tr.AccessToken)
	fmt.Println("refresh token: ", tr.RefreshToken)
	fmt.Println("id token: ", tr.IDToken)

	// Persist minimal identity fields by passwordless channel.
	if tr.IDToken != "" {
		if _, err := savePasswordlessIdentityFromIDToken(initializers.DB, tr.IDToken, c.Query("medium")); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	if c.Locals("isWebClient") == true {
		setAuthCookies(c, cfg, &tr)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Passwordless login confirmed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Passwordless login confirmed",
		"access_token":  tr.AccessToken,
		"refresh_token": tr.RefreshToken,
		"id_token":      tr.IDToken,
		"token_type":    tr.TokenType,
		"expires_in":    tr.ExpiresIn,
	})
}


// AuthCallback handles OAuth redirect from Auth0, exchanges the code, and sets httpOnly cookies.
func AuthCallback(c *fiber.Ctx) error {
	cfg := auth.GetAuthConfigs(c)
	if errMsg := c.Query("error"); errMsg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":             errMsg,
			"error_description": c.Query("error_description"),
		})
	}
	code := c.Query("code")
	stateQ := c.Query("state")
	stateC := c.Cookies(cfg.StateCookieName)
	clearStateCookie(c, cfg)

	if code == "" || stateQ == "" || stateC == "" ||
		subtle.ConstantTimeCompare([]byte(stateQ), []byte(stateC)) != 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid or missing OAuth state or code"})
	}

	// it's not returning refresh token. it returns only access token
	tr, err := auth.ExchangeAuthorizationCode(c.UserContext(), cfg, code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error1": err.Error()})
	}

	setAuthCookies(c, cfg, tr)

	if tr.IDToken != "" {
		if _, err := saveUserFromIDToken(initializers.DB, tr.IDToken); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.Redirect(cfg.PostLoginRedirect, fiber.StatusFound)
}

// AuthRefresh rotates the access token using the refresh token cookie.
// TODO: this also needs to know which configs to pick up, passwordless or normal ones
func AuthRefresh(c *fiber.Ctx) error {
	cfg := auth.Conf
	rt := strings.TrimSpace(c.Cookies(cfg.RefreshCookieName))
	if rt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing refresh token"})
	}
	tr, err := auth.RefreshTokens(c.UserContext(), cfg, rt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	setAuthCookies(c, cfg, tr)
	return c.JSON(fiber.Map{
		"token_type":         tr.TokenType,
		"expires_in":         tr.ExpiresIn,
		"refresh_expires_in": tr.RefreshExpiresIn,
	})
}

// AuthLogout clears app cookies and returns the Auth0 logout URL for the browser.
func AuthLogout(c *fiber.Ctx) error {
	cfg := auth.Conf
	clearAuthCookies(c, cfg)
	return c.JSON(fiber.Map{
		"logout_url": auth.LogoutURL(cfg),
	})
}

// AuthMe returns the authenticated user from the database.
func AuthMe(c *fiber.Ctx) error {
	sub, err := subjectFromAccessToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	var u models.User
	if err := initializers.DB.Where("auth0_id = ?", sub).First(&u).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(fiber.Map{"user": fiber.Map{
		"first_name": u.FirstName,
        "last_name": u.LastName,
        "email": u.Email,
        "phone_number": u.PhoneNumber,
        "email_verified": u.EmailVerified,
        "image_url": u.Image,
        "id": u.ID,
	}})
}

func subjectFromAccessToken(c *fiber.Ctx) (string, error) {
	token := auth.AccessTokenFromCtx(c)
	if token == "" {
		return "", fmt.Errorf("missing access token")
	}
	validated, _, err := auth.ValidateAccessTokenAny(c.UserContext(), token)
	if err != nil {
		return "", fmt.Errorf("invalid access token")
	}
	sub := strings.TrimSpace(validated.RegisteredClaims.Subject)
	if sub == "" {
		return "", fmt.Errorf("missing subject")
	}
	return sub, nil
}

func saveUserFromIDToken(db *gorm.DB, idToken string) (*models.User, error) {
	claims, err := auth.IDTokenClaims(idToken)
	if err != nil {
		return nil, err
	}
	u, err := userFromIDTokenClaims(claims)
	if err != nil {
		return nil, err
	}
	return upsertProfile(db, u.Auth0ID, models.ProfileInput{
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		Image:         u.Image,
		EmailVerified: u.EmailVerified,
	})
}

func savePasswordlessIdentityFromIDToken(db *gorm.DB, idToken string, medium string) (*models.User, error) {
	claims, err := auth.IDTokenClaims(idToken)
	if err != nil {
		return nil, err
	}
	sub := claimString(claims, "sub")
	if sub == "" {
		return nil, fmt.Errorf("id_token missing sub")
	}

	picture := claimString(claims, "picture")
	email_verified := claimBool(claims, "email_verified")

	email := ""
	phone := ""
	switch strings.ToLower(strings.TrimSpace(medium)) {
		case "sms":
			phone = claimString(claims, "name")
		default:
			email = claimString(claims, "email")
	}

	var existing models.User
	err = db.Where("auth0_id = ?", sub).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		row := models.User{
			Auth0ID:     sub,
			Email:       strings.TrimSpace(email),
			PhoneNumber: strings.TrimSpace(phone),
			Image:       picture,
			EmailVerified: email_verified,
		}
		if err := db.Create(&row).Error; err != nil {
			return nil, err
		}
		return &row, nil
	}
	if err != nil {
		return nil, err
	}

	if email != "" {
		existing.Email = strings.TrimSpace(email)
	}
	if phone != "" {
		existing.PhoneNumber = strings.TrimSpace(phone)
	}
	if picture != "" {
		existing.Image = picture
	}
	if err := db.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

// this will only work for the email redirection universal flow
// for self managed, all of this has to be input by user
func userFromIDTokenClaims(claims map[string]interface{}) (models.User, error) {
	sub := claimString(claims, "sub")
	if sub == "" {
		return models.User{}, fmt.Errorf("id_token missing sub")
	}

	first := claimString(claims, "given_name")
	last := claimString(claims, "family_name")
	if first == "" {
		first, last = splitName(claimString(claims, "name"))
	}

	email := claimString(claims, "email")
	if email == "" {
		return models.User{}, fmt.Errorf("id_token missing email")
	}

	return models.User{
		FirstName:     first,
		LastName:      last,
		Email:         email,
		PhoneNumber:   claimString(claims, "phone_number"),
		EmailVerified: claimBool(claims, "email_verified"),
		Image:         claimString(claims, "picture"),
		Auth0ID:       sub,
	}, nil
}

func claimString(claims map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		v, ok := claims[key]
		if !ok || v == nil {
			continue
		}
		if s, ok := v.(string); ok {
			if s = strings.TrimSpace(s); s != "" {
				return s
			}
		}
	}
	return ""
}

func claimBool(claims map[string]interface{}, key string) bool {
	v, ok := claims[key]
	if !ok || v == nil {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}

func splitName(full string) (first, last string) {
	full = strings.TrimSpace(full)
	if full == "" {
		return "", ""
	}
	parts := strings.SplitN(full, " ", 2)
	first = parts[0]
	if len(parts) > 1 {
		last = strings.TrimSpace(parts[1])
	}
	return first, last
}
