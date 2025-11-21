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

type problemFile struct {
	Contest string
	Index   string
	Path    string
}

func main() {
	root := flag.String("root", ".", "root directory that contains contest folders")
	dsn := flag.String("dsn", "", "database connection string (defaults to $DB_DSN or local postgres)")
	dryRun := flag.Bool("dry-run", false, "log planned inserts without touching the database")
	limit := flag.Int("limit", 0, "maximum number of inserts to perform (0 = no limit)")
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

	existing, err := loadExistingProblems(ctx, db)
	if err != nil {
		log.Fatalf("failed to load existing problems: %v", err)
	}

	idAlloc, err := newIDAllocator(ctx, db)
	if err != nil {
		log.Fatalf("failed to initialize id allocator: %v", err)
	}

	candidates, err := collectProblems(*root)
	if err != nil {
		log.Fatalf("failed to scan repository: %v", err)
	}

	if len(candidates) == 0 {
		log.Println("No problem statements found on disk.")
		return
	}

	var toInsert []problemFile
	for _, cand := range candidates {
		key := compositeKey(cand.Contest, cand.Index)
		if _, ok := existing[key]; ok {
			continue
		}
		toInsert = append(toInsert, cand)
	}

	if len(toInsert) == 0 {
		log.Println("All local problems already exist in the database.")
		return
	}

	log.Printf("Found %d new problems to insert.", len(toInsert))
	if *dryRun {
		for _, p := range toInsert {
			fmt.Printf("[dry-run] would insert contest %s problem %s from %s\n", p.Contest, p.Index, p.Path)
		}
		return
	}

	inserted := 0
	for _, p := range toInsert {
		if *limit > 0 && inserted >= *limit {
			break
		}
		if err := insertProblem(ctx, db, p, idAlloc); err != nil {
			log.Printf("failed inserting %s %s: %v", p.Contest, p.Index, err)
			continue
		}
		inserted++
	}

	log.Printf("Inserted %d problems (skipped %d already present).", inserted, len(toInsert)-inserted)
}

func collectProblems(root string) ([]problemFile, error) {
	var list []problemFile
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
		list = append(list, problemFile{
			Contest: contest,
			Index:   index,
			Path:    path,
		})
		return nil
	})
	return list, err
}

func compositeKey(contest, index string) string {
	return contest + "|" + strings.ToUpper(index)
}

func loadExistingProblems(ctx context.Context, db *sql.DB) (map[string]struct{}, error) {
	rows, err := db.QueryContext(ctx, "SELECT contest_id::text, UPPER(index_name) FROM problems")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existing := make(map[string]struct{})
	for rows.Next() {
		var contest, index string
		if err := rows.Scan(&contest, &index); err != nil {
			return nil, err
		}
		existing[compositeKey(strings.TrimSpace(contest), strings.TrimSpace(index))] = struct{}{}
	}
	return existing, rows.Err()
}

func insertProblem(ctx context.Context, db *sql.DB, p problemFile, alloc *idAllocator) error {
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return fmt.Errorf("read statement: %w", err)
	}
	stmtCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	title := deriveTitle(p, data)
	if alloc.needsManual {
		id := alloc.next()
		_, err = db.ExecContext(
			stmtCtx,
			`INSERT INTO problems (id, contest_id, index_name, title, statement) VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT DO NOTHING`,
			id,
			p.Contest,
			p.Index,
			title,
			string(data),
		)
	} else {
		_, err = db.ExecContext(
			stmtCtx,
			`INSERT INTO problems (contest_id, index_name, title, statement) VALUES ($1, $2, $3, $4)
			 ON CONFLICT DO NOTHING`,
			p.Contest,
			p.Index,
			title,
			string(data),
		)
	}
	return err
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
