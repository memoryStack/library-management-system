# Graph Report - library-management-system  (2026-05-15)

## Corpus Check
- 18 files · ~4,588 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 69 nodes · 88 edges · 10 communities detected
- Extraction: 76% EXTRACTED · 24% INFERRED · 0% AMBIGUOUS · INFERRED: 21 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]

## God Nodes (most connected - your core abstractions)
1. `AuthCallback()` - 6 edges
2. `init()` - 5 edges
3. `Init()` - 5 edges
4. `cookieSameSite()` - 5 edges
5. `doTokenForm()` - 4 edges
6. `RefreshTokens()` - 4 edges
7. `RequireAuth()` - 4 edges
8. `setAuthCookies()` - 4 edges
9. `AuthLogin()` - 4 edges
10. `MergeAuthCookiesIntoRequest()` - 4 edges

## Surprising Connections (you probably didn't know these)
- `init()` --calls--> `LoadEnv()`  [INFERRED]
  main.go → initializers/env.go
- `init()` --calls--> `ConnectDB()`  [INFERRED]
  main.go → initializers/db.go
- `init()` --calls--> `SyncDB()`  [INFERRED]
  main.go → initializers/db.go
- `init()` --calls--> `ConnectDB()`  [INFERRED]
  main.go → initializers/database.go
- `AuthCallback()` --calls--> `IDTokenClaims()`  [INFERRED]
  controllers/auth.go → auth/jwt.go

## Communities (13 total, 3 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.32
Nodes (11): GetAuthConfigs(), AuthorizeURL(), AuthCallback(), AuthLogin(), AuthLogout(), AuthRefresh(), clearAuthCookies(), clearStateCookie() (+3 more)

### Community 1 - "Community 1"
Cohesion: 0.22
Nodes (6): AccessTokenFromCtx(), IDTokenClaims(), initJWT(), ValidateAccessToken(), AuthMe(), RequireAuth()

### Community 2 - "Community 2"
Cohesion: 0.29
Nodes (5): init(), ConnectDB(), ConnectDB(), SyncDB(), LoadEnv()

### Community 3 - "Community 3"
Cohesion: 0.43
Nodes (5): Config, loadAuth0Config(), loadConfig(), loadPasswordlessConfigOptional(), Init()

### Community 4 - "Community 4"
Cohesion: 0.43
Nodes (6): doTokenForm(), ExchangeAuthorizationCode(), LogoutURL(), RefreshTokens(), tokenEndpoint(), TokenResponse

### Community 5 - "Community 5"
Cohesion: 0.53
Nodes (5): cookieSameSiteFiber(), joinCookieHeader(), MergeAuthCookiesIntoRequest(), parseCookieHeader(), SetAuthCookies()

### Community 6 - "Community 6"
Cohesion: 0.5
Nodes (3): main(), devCORSOrigins(), Stack()

## Knowledge Gaps
- **5 isolated node(s):** `Config`, `TokenResponse`, `User`, `Author`, `Book`
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `RequireAuth()` connect `Community 1` to `Community 4`, `Community 5`?**
  _High betweenness centrality (0.114) - this node is a cross-community bridge._
- **Why does `RefreshTokens()` connect `Community 4` to `Community 0`, `Community 1`?**
  _High betweenness centrality (0.085) - this node is a cross-community bridge._
- **Why does `MergeAuthCookiesIntoRequest()` connect `Community 5` to `Community 1`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `AuthCallback()` (e.g. with `GetAuthConfigs()` and `ExchangeAuthorizationCode()`) actually correct?**
  _`AuthCallback()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **Are the 4 inferred relationships involving `init()` (e.g. with `LoadEnv()` and `ConnectDB()`) actually correct?**
  _`init()` has 4 INFERRED edges - model-reasoned connections that need verification._
- **Are the 4 inferred relationships involving `Init()` (e.g. with `loadAuth0Config()` and `loadPasswordlessConfigOptional()`) actually correct?**
  _`Init()` has 4 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Config`, `TokenResponse`, `User` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._