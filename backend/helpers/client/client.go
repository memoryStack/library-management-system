package client

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ua-parser/uap-go/uaparser"
)

// Kind distinguishes browser (cookie session) from native (tokens in response body).
type Kind int

const (
	Web Kind = iota
	Native
)

// Info is the resolved client type for the current request.
type Info struct {
	Kind   Kind
	Source string // query, header, sec-fetch, ua-marker, ua-parser, or unknown
	OK     bool
}

// Resolve hybrid client detection: explicit query/header first, then Sec-Fetch,
// then raw UA markers, then uaparser as fallback.
//
// Accepted explicit values:
//   - web: web, web-app, browser, spa
//   - native: native, native-app, mobile, ios, android, app
func Resolve(c *fiber.Ctx, parser *uaparser.Parser) Info {
	if kind, ok := parseExplicit(c.Query("client")); ok {
		return Info{Kind: kind, Source: "query", OK: true}
	}
	if kind, ok := parseExplicit(c.Get("X-Client-Type")); ok {
		return Info{Kind: kind, Source: "header", OK: true}
	}
	if hasSecFetchMetadata(c) {
		return Info{Kind: Web, Source: "sec-fetch", OK: true}
	}

	ua := strings.TrimSpace(c.Get("User-Agent"))
	if ua == "" {
		return Info{Source: "unknown"}
	}

	if kind, ok := classifyRawUA(ua); ok {
		return Info{Kind: kind, Source: "ua-marker", OK: true}
	}

	if parser != nil {
		if kind, ok := classifyParsedUA(parser.Parse(ua)); ok {
			return Info{Kind: kind, Source: "ua-parser", OK: true}
		}
	}

	return Info{Source: "unknown"}
}

func parseExplicit(raw string) (Kind, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "web", "web-app", "browser", "spa":
		return Web, true
	case "native", "native-app", "mobile", "ios", "android", "app":
		return Native, true
	default:
		return 0, false
	}
}

func hasSecFetchMetadata(c *fiber.Ctx) bool {
	return c.Get("Sec-Fetch-Mode") != "" ||
		c.Get("Sec-Fetch-Site") != "" ||
		c.Get("Sec-Fetch-Dest") != ""
}

var nativeUAMarkers = []string{
	"okhttp", "alamofire", "cfnetwork", "darwin/",
	"reactnative", "react-native", "flutter", "dart/",
}

func classifyRawUA(ua string) (Kind, bool) {
	lower := strings.ToLower(ua)
	for _, m := range nativeUAMarkers {
		if strings.Contains(lower, m) {
			return Native, true
		}
	}
	return 0, false
}

var browserFamilies = map[string]struct{}{
	"Chrome":          {},
	"Chrome Mobile":   {},
	"Chromium":        {},
	"Firefox":         {},
	"Safari":          {},
	"Mobile Safari":   {},
	"Edge":            {},
	"Opera":           {},
	"Opera Mini":      {},
	"Samsung Internet": {},
	"IE":              {},
}

func classifyParsedUA(parsed *uaparser.Client) (Kind, bool) {
	if parsed == nil {
		return 0, false
	}
	family := parsed.UserAgent.Family
	if _, ok := browserFamilies[family]; ok {
		return Web, true
	}

	os := parsed.Os.Family
	if family == "Other" && (os == "Android" || os == "iOS" || os == "Windows Phone") {
		return Native, true
	}

	return 0, false
}

func (i Info) IsWeb() bool    { return i.OK && i.Kind == Web }
func (i Info) IsNative() bool { return i.OK && i.Kind == Native }
