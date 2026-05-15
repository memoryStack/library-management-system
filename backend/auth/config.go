package auth

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Config holds Auth0 OAuth/OIDC settings loaded from the environment.
// The same shape is used for the primary (e.g. Universal Login) app and for a parallel
// passwordless application, loaded with different env prefixes.
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

// Conf is the primary Auth0 application (AUTH0_*). Populated by Init.
var Conf *Config

// PasswordlessConf is the parallel passwordless Auth0 application (AUTH0_PASSWORDLESS_*).
// Nil when AUTH0_PASSWORDLESS_DOMAIN is unset (passwordless flow disabled).
var PasswordlessConf *Config

// loadAuth0Config loads one Auth0 app from env keys prefix + field suffix.
// Primary app uses prefix "AUTH0_" (e.g. AUTH0_DOMAIN). Passwordless uses "AUTH0_PASSWORDLESS_".
func loadAuth0Config(prefix string) (*Config, error) {
	key := func(suffix string) string {
		return prefix + suffix
	}
	get := func(suffix string) (string, error) {
		k := key(suffix)
		v := strings.TrimSpace(os.Getenv(k))
		if v == "" {
			return "", fmt.Errorf("missing required env %s", k)
		}
		return v, nil
	}

	domain, err := get("DOMAIN")
	if err != nil {
		return nil, err
	}
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimSuffix(domain, "/")

	clientID, err := get("CLIENT_ID")
	if err != nil {
		return nil, err
	}
	clientSecret, err := get("CLIENT_SECRET")
	if err != nil {
		return nil, err
	}
	callbackURL, err := get("CALLBACK_URL")
	if err != nil {
		return nil, err
	}
	audience, err := get("AUDIENCE")
	if err != nil {
		return nil, err
	}

	postLogin := strings.TrimSpace(os.Getenv(key("POST_LOGIN_REDIRECT")))
	if postLogin == "" {
		postLogin = "http://localhost:5173/"
	}
	logoutReturn := strings.TrimSpace(os.Getenv(key("LOGOUT_RETURN_URL")))
	if logoutReturn == "" {
		logoutReturn = postLogin
	}

	connection := strings.TrimSpace(os.Getenv(key("CONNECTION")))

	cookieSecure := false
	if v := strings.TrimSpace(os.Getenv(key("COOKIE_SECURE"))); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key("COOKIE_SECURE"), err)
		}
		cookieSecure = b
	}

	sameSite := strings.ToLower(strings.TrimSpace(os.Getenv(key("COOKIE_SAMESITE"))))
	if sameSite == "" {
		sameSite = "lax"
	}
	if sameSite == "none" && !cookieSecure {
		return nil, fmt.Errorf("%s=none requires %s=true", key("COOKIE_SAMESITE"), key("COOKIE_SECURE"))
	}

	accessName := strings.TrimSpace(os.Getenv(key("ACCESS_COOKIE_NAME")))
	if accessName == "" {
		if prefix == "AUTH0_" {
			accessName = "access_token"
		} else {
			accessName = "access_token_pw"
		}
	}
	refreshName := strings.TrimSpace(os.Getenv(key("REFRESH_COOKIE_NAME")))
	if refreshName == "" {
		if prefix == "AUTH0_" {
			refreshName = "refresh_token"
		} else {
			refreshName = "refresh_token_pw"
		}
	}
	stateName := strings.TrimSpace(os.Getenv(key("STATE_COOKIE_NAME")))
	if stateName == "" {
		stateName = "oauth_state"
	}

	refreshPath := strings.TrimSpace(os.Getenv(key("REFRESH_COOKIE_PATH")))
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

// loadPasswordlessConfigOptional returns nil if passwordless is not configured
// (AUTH0_PASSWORDLESS_DOMAIN unset); otherwise loads the full AUTH0_PASSWORDLESS_* set.
func loadPasswordlessConfigOptional() (*Config, error) {
	if strings.TrimSpace(os.Getenv("AUTH0_PASSWORDLESS_DOMAIN")) == "" {
		return nil, nil
	}
	return loadAuth0Config("AUTH0_PASSWORDLESS_")
}

func GetAuthConfigs(c *fiber.Ctx) (*Config) {
	medium := c.Query("medium")
	if medium == "" {
		return Conf
	} else if medium == "sms" {
		PasswordlessConf.Connection = "sms"
		return PasswordlessConf
	} else {
		PasswordlessConf.Connection = "email"
		return PasswordlessConf
	}
}
