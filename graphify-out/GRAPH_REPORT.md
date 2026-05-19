# Graph Report - library-management-system  (2026-05-20)

## Corpus Check
- 22 files · ~5,522 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 98 nodes · 126 edges · 11 communities detected
- Extraction: 76% EXTRACTED · 24% INFERRED · 0% AMBIGUOUS · INFERRED: 30 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]

## God Nodes (most connected - your core abstractions)
1. `DetectAuthClient()` - 6 edges
2. `AuthCallback()` - 6 edges
3. `Resolve()` - 6 edges
4. `resolveAuthClient()` - 6 edges
5. `init()` - 5 edges
6. `Init()` - 5 edges
7. `cookieSameSite()` - 5 edges
8. `setAuthCookies()` - 5 edges
9. `doTokenForm()` - 4 edges
10. `RefreshTokens()` - 4 edges

## Surprising Connections (you probably didn't know these)
- `init()` --calls--> `LoadEnv()`  [INFERRED]
  main.go → initializers/env.go
- `init()` --calls--> `ConnectDB()`  [INFERRED]
  main.go → initializers/db.go
- `init()` --calls--> `SyncDB()`  [INFERRED]
  main.go → initializers/db.go
- `init()` --calls--> `ConnectDB()`  [INFERRED]
  main.go → initializers/database.go
- `AuthCallback()` --calls--> `ExchangeAuthorizationCode()`  [INFERRED]
  controllers/auth.go → auth/oauth.go

## Communities (14 total, 3 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.21
Nodes (14): GetAuthConfigs(), AccessTokenFromCtx(), IDTokenClaims(), AuthCallback(), AuthConfirmOTP(), AuthLogin(), AuthLogout(), AuthMe() (+6 more)

### Community 1 - "Community 1"
Cohesion: 0.18
Nodes (8): init(), main(), ConnectDB(), ConnectDB(), SyncDB(), LoadEnv(), devCORSOrigins(), Stack()

### Community 2 - "Community 2"
Cohesion: 0.24
Nodes (8): cookieSameSiteFiber(), joinCookieHeader(), MergeAuthCookiesIntoRequest(), parseCookieHeader(), SetAuthCookies(), initJWT(), ValidateAccessToken(), RequireAuth()

### Community 3 - "Community 3"
Cohesion: 0.29
Nodes (7): classifyParsedUA(), classifyRawUA(), hasSecFetchMetadata(), parseExplicit(), Resolve(), Info, Kind

### Community 4 - "Community 4"
Cohesion: 0.22
Nodes (5): Attach(), FromContext(), SetLocals(), Parser(), DetectAuthClient()

### Community 5 - "Community 5"
Cohesion: 0.36
Nodes (7): AuthorizeURL(), doTokenForm(), ExchangeAuthorizationCode(), LogoutURL(), RefreshTokens(), tokenEndpoint(), TokenResponse

### Community 6 - "Community 6"
Cohesion: 0.36
Nodes (6): detectAuthClientFromRequest(), parseAuthClientValue(), resolveAuthClient(), TestResolveAuthClient_explicit(), TestResolveAuthClient_secFetch(), AuthClientKind

### Community 7 - "Community 7"
Cohesion: 0.43
Nodes (5): Config, loadAuth0Config(), loadConfig(), loadPasswordlessConfigOptional(), Init()

## Knowledge Gaps
- **7 isolated node(s):** `Config`, `TokenResponse`, `User`, `Author`, `Book` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `AuthConfirmOTP()` connect `Community 0` to `Community 4`, `Community 6`?**
  _High betweenness centrality (0.416) - this node is a cross-community bridge._
- **Why does `DetectAuthClient()` connect `Community 4` to `Community 1`, `Community 3`?**
  _High betweenness centrality (0.379) - this node is a cross-community bridge._
- **Why does `FromContext()` connect `Community 4` to `Community 0`?**
  _High betweenness centrality (0.340) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `DetectAuthClient()` (e.g. with `Parser()` and `Resolve()`) actually correct?**
  _`DetectAuthClient()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 3 inferred relationships involving `AuthCallback()` (e.g. with `GetAuthConfigs()` and `ExchangeAuthorizationCode()`) actually correct?**
  _`AuthCallback()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **Are the 3 inferred relationships involving `resolveAuthClient()` (e.g. with `AuthConfirmOTP()` and `TestResolveAuthClient_explicit()`) actually correct?**
  _`resolveAuthClient()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **Are the 4 inferred relationships involving `init()` (e.g. with `LoadEnv()` and `ConnectDB()`) actually correct?**
  _`init()` has 4 INFERRED edges - model-reasoned connections that need verification._