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

	// "github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gofiber/fiber/v2"

	"library-management-system/backend/initializers"
	"library-management-system/backend/models"
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
	defer resp.Body.Close()

	if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Passwordless login started",
	})
}

type OTPConfirmResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken string `json:"id_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
	Scope string `json:"scope"`
}

func AuthConfirmOTP(c *fiber.Ctx) error {
	cfg := auth.PasswordlessConf
	
	url := "https://"+auth.PasswordlessConf.Domain+"/oauth/token"
	
	var requestBody map[string]string
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	requestPayload := map[string]string{
		"grant_type": "http://auth0.com/oauth/grant-type/passwordless/otp",
		"client_id": cfg.ClientID,
		"client_secret": cfg.ClientSecret,
		"username": requestBody["username"],
		"otp": requestBody["otp"],
		"realm": c.Query("medium"),
		"scope": "openid profile email offline_access",
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	defer resp.Body.Close()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var response OTPConfirmResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Passwordless login confirmed",
		"access_token": response.AccessToken,
		"refresh_token": response.RefreshToken,
		"id_token": response.IDToken,
		"token_type": response.TokenType,
		"expires_in": response.ExpiresIn,
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

	// save user to database
	idtoken := tr.IDToken
	userValues, err := auth.IDTokenClaims(idtoken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error2": err.Error()})
	}

	present := initializers.DB.Where("auth0_id = ?", userValues["sub"].(string)).First(&models.User{})
	if present.Error == nil {
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user already exists"})
		return c.Redirect(cfg.PostLoginRedirect, fiber.StatusFound)
	}

	result := initializers.DB.Create(&models.User{
		Name: userValues["name"].(string),
		Email: userValues["email"].(string),
		EmailVerified: userValues["email_verified"].(bool),
		Image: userValues["picture"].(string),
		Auth0ID: userValues["sub"].(string),
	})

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "some error while saving user details"})
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

// AuthMe returns the authenticated subject (JWT already validated by middleware).
func AuthMe(c *fiber.Ctx) error {
	idtoken := auth.AccessTokenFromCtx(c)
	if idtoken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing access token"})
	}
	userValues, err := auth.IDTokenClaims(idtoken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	user := models.User{}
	result := initializers.DB.Where("auth0_id = ?", userValues["sub"].(string)).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "some error while getting user details"})
	}
	return c.JSON(fiber.Map{
		"user": user,
	})
}
