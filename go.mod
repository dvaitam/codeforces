module codeforces

// Set to a broadly compatible version to avoid toolchain
// mismatches on environments older than Go 1.21.
go 1.20

// Pin MySQL driver to a version that supports Go <= 1.20.
require github.com/go-sql-driver/mysql v1.7.1

require filippo.io/edwards25519 v1.1.0 // indirect
