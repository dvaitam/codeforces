package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"
)

// cfAPIResponse mirrors the top-level Codeforces API envelope for problemset.problems.
type cfAPIResponse struct {
	Status  string          `json:"status"`
	Result  cfProblemResult `json:"result"`
	Comment string          `json:"comment"`
}

type cfProblemResult struct {
	Problems []cfProblem `json:"problems"`
}

// cfStandingsResponse mirrors the envelope for contest.standings.
type cfStandingsResponse struct {
	Status  string         `json:"status"`
	Result  cfStandingsResult `json:"result"`
	Comment string         `json:"comment"`
}

type cfStandingsResult struct {
	Problems []cfProblem `json:"problems"`
}

type cfProblem struct {
	ContestID int      `json:"contestId"`
	Index     string   `json:"index"`
	Name      string   `json:"name"`
	Rating    int      `json:"rating"`
	Tags      []string `json:"tags"`
}

const cfProblemsURL = "https://codeforces.com/api/problemset.problems"

var httpClient = &http.Client{Timeout: 60 * time.Second}

func main() {
	dsn := flag.String("dsn", "", "Postgres DSN (defaults to $DB_DSN or local postgres)")
	dryRun := flag.Bool("dry-run", false, "print planned updates without touching the database")
	insertMissing := flag.Bool("insert-missing", false, "insert problems from the API that are not yet in the database")
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
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	pingCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		log.Fatalf("connect db: %v", err)
	}

	ctx := context.Background()

	if err := ensureMetadataColumns(ctx, db); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}

	// --- Pass 1: problemset.problems ---
	log.Printf("Fetching problem list from %s …", cfProblemsURL)
	problems, err := fetchCFProblems()
	if err != nil {
		log.Fatalf("fetch CF API: %v", err)
	}
	log.Printf("Fetched %d problems from Codeforces API.", len(problems))

	// Track which contest IDs the problemset API covered.
	coveredContests := make(map[string]bool)
	for _, p := range problems {
		coveredContests[fmt.Sprintf("%d", p.ContestID)] = true
	}

	updated, inserted, skipped := applyProblems(ctx, db, problems, *dryRun, *insertMissing)
	if !*dryRun {
		log.Printf("Pass 1 done: %d updated, %d inserted, %d not found.", updated, inserted, skipped)
	}

	// --- Pass 2: contest.standings for contests the problemset API missed ---
	dbContests, err := loadDBContestIDs(ctx, db)
	if err != nil {
		log.Fatalf("load contest ids: %v", err)
	}

	var missed []string
	for _, c := range dbContests {
		if !coveredContests[c] {
			missed = append(missed, c)
		}
	}

	if len(missed) == 0 {
		log.Println("All DB contests were covered by problemset.problems — nothing more to do.")
		return
	}

	log.Printf("Pass 2: %d contest(s) not in problemset.problems, fetching via contest.standings …", len(missed))
	u2, i2, s2 := 0, 0, 0
	for idx, contestID := range missed {
		// Respect the CF API rate limit (~1 req/s for unauthenticated).
		if idx > 0 {
			time.Sleep(time.Second)
		}
		probs, err := fetchContestProblems(contestID)
		if err != nil {
			log.Printf("standings %s: %v", contestID, err)
			continue
		}
		u, i, s := applyProblems(ctx, db, probs, *dryRun, *insertMissing)
		u2 += u
		i2 += i
		s2 += s
	}

	if !*dryRun {
		log.Printf("Pass 2 done: %d updated, %d inserted, %d not found.", u2, i2, s2)
	}
}

// applyProblems updates (and optionally inserts) a slice of CF problems in the DB.
func applyProblems(ctx context.Context, db *sql.DB, problems []cfProblem, dryRun, insertMissing bool) (updated, inserted, skipped int) {
	for _, p := range problems {
		contestID := fmt.Sprintf("%d", p.ContestID)
		if dryRun {
			log.Printf("[dry-run] contest=%s index=%s title=%q rating=%d tags=%s",
				contestID, p.Index, p.Name, p.Rating, strings.Join(p.Tags, ","))
			continue
		}
		n, err := updateProblemMetadata(ctx, db, contestID, p.Index, p.Name, p.Rating, p.Tags)
		if err != nil {
			log.Printf("update %s%s: %v", contestID, p.Index, err)
			continue
		}
		if n > 0 {
			updated++
			continue
		}
		if insertMissing {
			if err := insertProblemFromAPI(ctx, db, contestID, p.Index, p.Name, p.Rating, p.Tags); err != nil {
				log.Printf("insert %s%s: %v", contestID, p.Index, err)
				continue
			}
			inserted++
		} else {
			skipped++
		}
	}
	return
}

// loadDBContestIDs returns all distinct contest_id values present in the DB.
func loadDBContestIDs(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT DISTINCT contest_id FROM problems ORDER BY contest_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, strings.TrimSpace(id))
	}
	return ids, rows.Err()
}

// ensureMetadataColumns adds rating and tags columns if they don't exist yet.
func ensureMetadataColumns(ctx context.Context, db *sql.DB) error {
	stmts := []string{
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS rating INT`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS tags TEXT[]`,
	}
	for _, s := range stmts {
		if _, err := db.ExecContext(ctx, s); err != nil {
			return fmt.Errorf("%s: %w", s, err)
		}
	}
	return nil
}

// fetchCFProblems fetches the full problemset from the Codeforces API.
func fetchCFProblems() ([]cfProblem, error) {
	resp, err := httpClient.Get(cfProblemsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var apiResp cfAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("API status %q: %s", apiResp.Status, apiResp.Comment)
	}
	return apiResp.Result.Problems, nil
}

// fetchContestProblems fetches the problem list for a single contest via standings.
func fetchContestProblems(contestID string) ([]cfProblem, error) {
	url := fmt.Sprintf("https://codeforces.com/api/contest.standings?contestId=%s&from=1&count=1", contestID)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var apiResp cfStandingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("API status %q: %s", apiResp.Status, apiResp.Comment)
	}
	return apiResp.Result.Problems, nil
}

// insertProblemFromAPI inserts a stub row sourced from the CF API.
// statement, reference_solution, and verifier are left NULL.
func insertProblemFromAPI(ctx context.Context, db *sql.DB, contestID, index, title string, rating int, tags []string) error {
	var ratingVal interface{}
	if rating > 0 {
		ratingVal = rating
	}
	_, err := db.ExecContext(ctx,
		`INSERT INTO problems (contest_id, index_name, title, rating, tags)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT DO NOTHING`,
		contestID, index, title, ratingVal, pq.Array(tags),
	)
	return err
}

// updateProblemMetadata updates title, rating, and tags for a matching row.
// Returns the number of rows affected (0 if no row matched).
func updateProblemMetadata(ctx context.Context, db *sql.DB, contestID, index, title string, rating int, tags []string) (int64, error) {
	var ratingVal interface{}
	if rating > 0 {
		ratingVal = rating
	}
	res, err := db.ExecContext(ctx,
		`UPDATE problems
		    SET title  = $1,
		        rating = $2,
		        tags   = $3
		  WHERE contest_id = $4
		    AND UPPER(index_name) = UPPER($5)`,
		title, ratingVal, pq.Array(tags), contestID, index,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
