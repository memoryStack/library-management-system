module library-management-system/backend

go 1.24.0

require (
	github.com/auth0/go-jwt-middleware/v2 v2.3.1
	github.com/gofiber/fiber/v2 v2.52.5
	github.com/joho/godotenv v1.5.1
	github.com/ua-parser/uap-go v0.0.0-20251207011819-db9adb27a0b8
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.12
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	gopkg.in/go-jose/go-jose.v2 v2.6.3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Go 1.26+ can mis-resolve go-spew v1.1.1 in some transitive test graphs; pin to a commit that contains ./spew.
replace github.com/davecgh/go-spew v1.1.1 => github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc

// Same class of issue for go-difflib v1.0.0 tags vs package layout.
replace github.com/pmezard/go-difflib v1.0.0 => github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2

// v3.0.1 zip/module metadata can fail "does not contain package" on some toolchains; v3.0.0 resolves cleanly.
replace gopkg.in/yaml.v3 v3.0.1 => gopkg.in/yaml.v3 v3.0.0
