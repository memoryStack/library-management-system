# Project prompts

This file records the user prompts that shaped this repository, in conversation order.

---

## Prompt 1 — Go backend scaffold

Create a new Go project in a folder named `backend`. The backend is a Go server with these dependencies:

1. **go fiber** — routing  
2. **gorm** — database  
3. **godotenv** — environment variables  
4. **compiledaemon** — file watching and server restarts  

Required layout:

**Folders**

1. `controllers` — route handlers and business logic  
2. `helpers` — reusable helpers  
3. `initializers` — env loading, DB connection, etc., invoked from `main`  
4. `models` — data models  
5. `middlewares` — HTTP middlewares (with recommendations for a robust API)  

**Files**

1. `main.go` — routes, DB init, server start  

**Environments**

Two modes: **development** and **production**, with different DB (and related) setup. Load the matching env file when running the dev or prod server.

**Database:** PostgreSQL.

---

## Prompt 2 — Documentation and simpler configuration

1. Document all prompts for this project in a file in the repository.  
2. Remove “fallback-heavy” configuration (e.g. `BACKEND_ROOT`, `APP_ENV`, and similar optional paths). Prefer a **simple, conservative** setup with fewer configuration knobs.

---

## Prompt 3 — Committed env files and sharing with other developers

Because `initializers/env.go` loads `.env.development` or `.env.production`, create those files (not only `.env.{env}.example`) with everything the app needs.

Also: those filenames were listed in `backend/.gitignore`; the goal is to **share** them with other developers. Is committing them (instead of only examples) a good practice?

---

## How this maps to the code

- **Prompt 1** → `backend/` module (Fiber, GORM, godotenv, Makefile target for CompileDaemon, `controllers`, `helpers`, `initializers`, `models`, `middlewares`, `main.go`, `.env.development` / `.env.production` pattern).  
- **Prompt 2** → This file (`PROMPTS.md`) plus reduced env surface and stricter startup rules in `backend` (see comments in `main.go` and `initializers`).  
- **Prompt 3** → Tracked `backend/.env.development` and `backend/.env.production` with the required variables and comments; removed those paths from `backend/.gitignore`; dropped the separate `.example` files so there is one source of truth. Production file uses placeholders—real secrets should not be pushed (especially to a public repo).
