package main

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type problemAsset struct {
	Contest   string
	Index     string
	Statement string
	Title     string
	Solution  string
	Verifier  string
	Path      string
}

func main() {
	root := flag.String("root", ".", "root directory that contains contest folders")
	dsn := flag.String("dsn", "", "database connection string (defaults to $DB_DSN or local postgres)")
	dryRun := flag.Bool("dry-run", false, "log planned upserts without touching the database")
	limit := flag.Int("limit", 0, "maximum number of rows to write (0 = no limit)")
	contestFilter := flag.String("contest", "", "only sync this contest id")
	problemFilter := flag.String("problem", "", "only sync a single problem, e.g. 1994A")
	flag.Parse()

	connStr := *dsn
	if connStr == "" {
		connStr = os.Getenv("DB_DSN")
	}
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost:5432/codeforces?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	connectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(connectCtx); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	ctx := context.Background()

	if err := ensureProblemsSchema(ctx, db); err != nil {
		log.Fatalf("failed to ensure schema: %v", err)
	}

	idAlloc, err := newIDAllocator(ctx, db)
	if err != nil {
		log.Fatalf("failed to initialize id allocator: %v", err)
	}

	contestFilterVal := strings.TrimSpace(*contestFilter)
	problemContest, problemIndex, err := parseProblemArg(*problemFilter)
	if err != nil {
		log.Fatalf("invalid -problem value: %v", err)
	}
	if problemContest != "" && contestFilterVal != "" && problemContest != contestFilterVal {
		log.Fatalf("-problem contest %s does not match -contest %s", problemContest, contestFilterVal)
	}
	if contestFilterVal == "" {
		contestFilterVal = problemContest
	}

	assets, err := collectProblemAssets(*root, contestFilterVal, problemContest, problemIndex)
	if err != nil {
		log.Fatalf("failed to scan repository: %v", err)
	}

	if len(assets) == 0 {
		log.Println("No problem statements found on disk.")
		return
	}

	written := 0
	inserted := 0
	updated := 0

	for _, asset := range assets {
		if *limit > 0 && written >= *limit {
			break
		}
		if *dryRun {
			log.Printf("[dry-run] upsert %s %s from %s (solution: %t, verifier: %t)", asset.Contest, asset.Index, asset.Path, asset.Solution != "", asset.Verifier != "")
			written++
			continue
		}

		action, err := upsertProblem(ctx, db, asset, idAlloc)
		if err != nil {
			log.Printf("failed upserting %s %s: %v", asset.Contest, asset.Index, err)
			continue
		}
		written++
		if action == "inserted" {
			inserted++
		} else if action == "updated" {
			updated++
		}
	}

	log.Printf("Finished: %d inserted, %d updated, %d total processed.", inserted, updated, written)
}

func collectProblemAssets(root, contestFilter, problemContest, problemIndex string) ([]problemAsset, error) {
	var list []problemAsset
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if !strings.HasPrefix(name, "problem") || !strings.HasSuffix(name, ".txt") {
			return nil
		}
		index := strings.TrimSuffix(strings.TrimPrefix(name, "problem"), ".txt")
		if index == "" {
			return nil
		}
		dir := filepath.Dir(path)
		contest := filepath.Base(dir)
		if _, err := strconv.Atoi(contest); err != nil {
			return nil
		}
		if contestFilter != "" && contest != contestFilter {
			return nil
		}
		if problemIndex != "" && (contest != problemContest || strings.ToUpper(index) != problemIndex) {
			return nil
		}
		stmtData, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read statement %s: %w", path, err)
		}
		title := deriveTitle(problemFile{Contest: contest, Index: index}, stmtData)
		solutionPath := filepath.Join(dir, contest+index+".go")
		solutionData, _ := os.ReadFile(solutionPath)
		verifierPath := findVerifier(dir, index)
		verifierData, _ := os.ReadFile(verifierPath)
		list = append(list, problemAsset{
			Contest:   contest,
			Index:     strings.ToUpper(index),
			Statement: string(stmtData),
			Title:     title,
			Solution:  string(solutionData),
			Verifier:  string(verifierData),
			Path:      path,
		})
		return nil
	})
	return list, err
}

func parseProblemArg(val string) (string, string, error) {
	v := strings.TrimSpace(val)
	if v == "" {
		return "", "", nil
	}
	split := -1
	for i, r := range v {
		if r < '0' || r > '9' {
			split = i
			break
		}
	}
	if split <= 0 || split >= len(v) {
		return "", "", fmt.Errorf("expected contest digits followed by index, e.g. 1994A")
	}
	contest := v[:split]
	if _, err := strconv.Atoi(contest); err != nil {
		return "", "", fmt.Errorf("contest must be numeric: %w", err)
	}
	index := strings.ToUpper(strings.TrimSpace(v[split:]))
	if index == "" {
		return "", "", fmt.Errorf("missing problem index")
	}
	return contest, index, nil
}

func ensureProblemsSchema(ctx context.Context, db *sql.DB) error {
	ddl := []string{
		`CREATE TABLE IF NOT EXISTS problems (
			id SERIAL PRIMARY KEY,
			contest_id VARCHAR(20),
			index_name VARCHAR(20),
			title TEXT,
			statement TEXT,
			reference_solution TEXT,
			verifier TEXT
		)`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS title TEXT`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS statement TEXT`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS contest_id VARCHAR(20)`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS index_name VARCHAR(20)`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS reference_solution TEXT`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS verifier TEXT`,
	}
	for _, stmt := range ddl {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

func upsertProblem(ctx context.Context, db *sql.DB, asset problemAsset, alloc *idAllocator) (string, error) {
	var id int64
	err := db.QueryRowContext(ctx, `SELECT id FROM problems WHERE contest_id = $1 AND UPPER(index_name) = UPPER($2)`, asset.Contest, asset.Index).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == sql.ErrNoRows {
		stmtCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if alloc.needsManual {
			id = alloc.next()
			_, err = db.ExecContext(
				stmtCtx,
				`INSERT INTO problems (id, contest_id, index_name, title, statement, reference_solution, verifier) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				id,
				asset.Contest,
				asset.Index,
				asset.Title,
				asset.Statement,
				asset.Solution,
				asset.Verifier,
			)
		} else {
			_, err = db.ExecContext(
				stmtCtx,
				`INSERT INTO problems (contest_id, index_name, title, statement, reference_solution, verifier) VALUES ($1, $2, $3, $4, $5, $6)`,
				asset.Contest,
				asset.Index,
				asset.Title,
				asset.Statement,
				asset.Solution,
				asset.Verifier,
			)
		}
		return "inserted", err
	}

	stmtCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = db.ExecContext(
		stmtCtx,
		`UPDATE problems SET title = $1, statement = $2, reference_solution = $3, verifier = $4 WHERE id = $5`,
		asset.Title,
		asset.Statement,
		asset.Solution,
		asset.Verifier,
		id,
	)
	return "updated", err
}

func deriveTitle(p problemFile, data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			return line
		}
	}
	return fmt.Sprintf("Contest %s Problem %s", p.Contest, p.Index)
}

type idAllocator struct {
	needsManual bool
	current     int64
}

func newIDAllocator(ctx context.Context, db *sql.DB) (*idAllocator, error) {
	var colDefault sql.NullString
	err := db.QueryRowContext(ctx, `
		SELECT column_default
		FROM information_schema.columns
		WHERE table_name = 'problems' AND column_name = 'id'
	`).Scan(&colDefault)
	if err != nil {
		return nil, err
	}
	if colDefault.Valid && strings.Contains(colDefault.String, "nextval") {
		return &idAllocator{needsManual: false}, nil
	}
	var maxID sql.NullInt64
	if err := db.QueryRowContext(ctx, `SELECT COALESCE(MAX(id), 0) FROM problems`).Scan(&maxID); err != nil {
		return nil, err
	}
	return &idAllocator{
		needsManual: true,
		current:     maxID.Int64,
	}, nil
}

func (a *idAllocator) next() int64 {
	a.current++
	return a.current
}

func findVerifier(dir, letter string) string {
	cand := filepath.Join(dir, "verifier"+letter+".go")
	if _, err := os.Stat(cand); err == nil {
		return cand
	}
	cand = filepath.Join(dir, "verifier.go")
	if _, err := os.Stat(cand); err == nil {
		return cand
	}
	return ""
}

type problemFile struct {
	Contest string
	Index   string
	Path    string
}
