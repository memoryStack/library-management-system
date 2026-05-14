package auth

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds Auth0 OAuth/OIDC settings loaded from the environment.
type Config struct {
	Domain            string
	ClientID          string
	ClientSecret      string
	CallbackURL       string
	Audience          string // Auth0 API identifier — access JWTs must include this aud.
	PostLoginRedirect string // Browser redirect after successful callback (e.g. SPA URL).
	LogoutReturnURL   string // Allowed return URL for Auth0 /v2/logout (often same as PostLogin).

	Connection string // Optional Auth0 connection name (e.g. email passwordless).

	CookieSecure      bool
	CookieSameSite    string // lax | strict | none
	AccessCookieName  string
	RefreshCookieName string
	StateCookieName   string
	RefreshCookiePath string // narrow path so refresh token is not sent to every API route
}

// Conf is populated by Init and read by controllers/middlewares.
var Conf *Config

func loadConfig() (*Config, error) {
	get := func(key string) (string, error) {
		v := strings.TrimSpace(os.Getenv(key))
		if v == "" {
			return "", fmt.Errorf("missing required env %s", key)
		}
		return v, nil
	}

	domain, err := get("AUTH0_DOMAIN")
	if err != nil {
		return nil, err
	}
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimSuffix(domain, "/")

	clientID, err := get("AUTH0_CLIENT_ID")
	if err != nil {
		return nil, err
	}
	clientSecret, err := get("AUTH0_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}
	callbackURL, err := get("AUTH0_CALLBACK_URL")
	if err != nil {
		return nil, err
	}
	audience, err := get("AUTH0_AUDIENCE")
	if err != nil {
		return nil, err
	}

	postLogin := strings.TrimSpace(os.Getenv("AUTH0_POST_LOGIN_REDIRECT"))
	if postLogin == "" {
		postLogin = "http://localhost:5173/"
	}
	logoutReturn := strings.TrimSpace(os.Getenv("AUTH0_LOGOUT_RETURN_URL"))
	if logoutReturn == "" {
		logoutReturn = postLogin
	}

	connection := strings.TrimSpace(os.Getenv("AUTH0_CONNECTION"))

	cookieSecure := false
	if v := strings.TrimSpace(os.Getenv("AUTH0_COOKIE_SECURE")); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("AUTH0_COOKIE_SECURE: %w", err)
		}
		cookieSecure = b
	}

	sameSite := strings.ToLower(strings.TrimSpace(os.Getenv("AUTH0_COOKIE_SAMESITE")))
	if sameSite == "" {
		sameSite = "lax"
	}
	if sameSite == "none" && !cookieSecure {
		return nil, fmt.Errorf("AUTH0_COOKIE_SAMESITE=none requires AUTH0_COOKIE_SECURE=true")
	}

	accessName := strings.TrimSpace(os.Getenv("AUTH0_ACCESS_COOKIE_NAME"))
	if accessName == "" {
		accessName = "access_token"
	}
	refreshName := strings.TrimSpace(os.Getenv("AUTH0_REFRESH_COOKIE_NAME"))
	if refreshName == "" {
		refreshName = "refresh_token"
	}
	stateName := strings.TrimSpace(os.Getenv("AUTH0_STATE_COOKIE_NAME"))
	if stateName == "" {
		stateName = "oauth_state"
	}

	refreshPath := strings.TrimSpace(os.Getenv("AUTH0_REFRESH_COOKIE_PATH"))
	if refreshPath == "" {
		refreshPath = "/api/auth"
	}

	return &Config{
		Domain:            domain,
		ClientID:          clientID,
		ClientSecret:      clientSecret,
		CallbackURL:       callbackURL,
		Audience:          audience,
		PostLoginRedirect: postLogin,
		LogoutReturnURL:   logoutReturn,
		Connection:        connection,
		CookieSecure:      cookieSecure,
		CookieSameSite:    sameSite,
		AccessCookieName:  accessName,
		RefreshCookieName: refreshName,
		StateCookieName:   stateName,
		RefreshCookiePath: refreshPath,
	}, nil
}
