# Graph Report - library-management-system  (2026-05-20)

## Corpus Check
- 24 files · ~6,659 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 133 nodes · 200 edges · 16 communities detected
- Extraction: 77% EXTRACTED · 23% INFERRED · 0% AMBIGUOUS · INFERRED: 46 edges (avg confidence: 0.8)
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
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]

## God Nodes (most connected - your core abstractions)
1. `AuthCallback()` - 8 edges
2. `RequireAuth()` - 7 edges
3. `AuthConfirmOTP()` - 7 edges
4. `saveUserFromIDToken()` - 7 edges
5. `IDTokenClaims()` - 6 edges
6. `Init()` - 6 edges
7. `DetectAuthClient()` - 6 edges
8. `subjectFromAccessToken()` - 6 edges
9. `Resolve()` - 6 edges
10. `SaveFromIDToken()` - 6 edges

## Surprising Connections (you probably didn't know these)
- `AuthConfirmOTP()` --calls--> `FromContext()`  [INFERRED]
  controllers/auth.go → helpers/client/context.go
- `init()` --calls--> `LoadEnv()`  [INFERRED]
  main.go → initializers/env.go
- `init()` --calls--> `ConnectDB()`  [INFERRED]
  main.go → initializers/db.go
- `init()` --calls--> `SyncDB()`  [INFERRED]
  main.go → initializers/db.go
- `init()` --calls--> `ConnectDB()`  [INFERRED]
  main.go → initializers/database.go

## Communities (19 total, 4 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.18
Nodes (20): GetAuthConfigs(), IDTokenClaims(), AuthorizeURL(), LogoutURL(), AuthCallback(), AuthConfirmOTP(), AuthLogin(), AuthLogout() (+12 more)

### Community 1 - "Community 1"
Cohesion: 0.18
Nodes (8): init(), main(), ConnectDB(), ConnectDB(), SyncDB(), LoadEnv(), devCORSOrigins(), Stack()

### Community 2 - "Community 2"
Cohesion: 0.32
Nodes (10): AuthConfigs(), initJWT(), initJWTPasswordless(), newJWTValidator(), ValidateAccessToken(), ValidateAccessTokenAny(), ValidateAccessTokenForConfig(), validateWith() (+2 more)

### Community 3 - "Community 3"
Cohesion: 0.31
Nodes (8): AccessTokenFromCtx(), AuthMe(), subjectFromAccessToken(), isUserFieldMissing(), ProfileRequirements(), upsertProfile(), UpsertUserProfile(), validateBackupFieldsByMethod()

### Community 4 - "Community 4"
Cohesion: 0.29
Nodes (7): classifyParsedUA(), classifyRawUA(), hasSecFetchMetadata(), parseExplicit(), Resolve(), Info, Kind

### Community 5 - "Community 5"
Cohesion: 0.22
Nodes (5): Attach(), FromContext(), SetLocals(), Parser(), DetectAuthClient()

### Community 6 - "Community 6"
Cohesion: 0.43
Nodes (7): claimBool(), claimString(), SaveFromIDToken(), splitName(), UpsertProfile(), UserFromIDTokenClaims(), ProfileInput

### Community 7 - "Community 7"
Cohesion: 0.36
Nodes (6): detectAuthClientFromRequest(), parseAuthClientValue(), resolveAuthClient(), TestResolveAuthClient_explicit(), TestResolveAuthClient_secFetch(), AuthClientKind

### Community 8 - "Community 8"
Cohesion: 0.43
Nodes (5): Config, loadAuth0Config(), loadConfig(), loadPasswordlessConfigOptional(), Init()

### Community 9 - "Community 9"
Cohesion: 0.53
Nodes (5): doTokenForm(), ExchangeAuthorizationCode(), RefreshTokens(), tokenEndpoint(), TokenResponse

### Community 10 - "Community 10"
Cohesion: 0.53
Nodes (5): cookieSameSiteFiber(), joinCookieHeader(), MergeAuthCookiesIntoRequest(), parseCookieHeader(), SetAuthCookies()

### Community 12 - "Community 12"
Cohesion: 0.7
Nodes (4): ClaimBool(), ClaimString(), SplitName(), UserFromIDTokenClaims()

## Knowledge Gaps
- **9 isolated node(s):** `Config`, `TokenResponse`, `User`, `Author`, `ProfileInput` (+4 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `AuthConfirmOTP()` connect `Community 0` to `Community 5`, `Community 6`, `Community 7`?**
  _High betweenness centrality (0.386) - this node is a cross-community bridge._
- **Why does `DetectAuthClient()` connect `Community 5` to `Community 1`, `Community 4`?**
  _High betweenness centrality (0.311) - this node is a cross-community bridge._
- **Why does `FromContext()` connect `Community 5` to `Community 0`?**
  _High betweenness centrality (0.301) - this node is a cross-community bridge._
- **Are the 4 inferred relationships involving `AuthCallback()` (e.g. with `GetAuthConfigs()` and `ExchangeAuthorizationCode()`) actually correct?**
  _`AuthCallback()` has 4 INFERRED edges - model-reasoned connections that need verification._
- **Are the 6 inferred relationships involving `RequireAuth()` (e.g. with `ValidateAccessTokenAny()` and `AuthConfigs()`) actually correct?**
  _`RequireAuth()` has 6 INFERRED edges - model-reasoned connections that need verification._
- **Are the 3 inferred relationships involving `AuthConfirmOTP()` (e.g. with `SaveFromIDToken()` and `FromContext()`) actually correct?**
  _`AuthConfirmOTP()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **Are the 3 inferred relationships involving `saveUserFromIDToken()` (e.g. with `IDTokenClaims()` and `upsertProfile()`) actually correct?**
  _`saveUserFromIDToken()` has 3 INFERRED edges - model-reasoned connections that need verification._