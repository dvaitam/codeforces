package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type missingProblem struct {
	ContestID int
	Index     string
	Name      string
	URL       string
}

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

const (
	defaultUserAgent    = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	defaultFetchTimeout = 30 * time.Second
)

func main() {
	csvPath := flag.String("csv", "truly_missing_problems.csv", "path to truly_missing_problems.csv")
	dryRun := flag.Bool("dry-run", false, "log actions without fetching or writing files")
	start := flag.Int("start", 0, "skip the first N entries")
	limit := flag.Int("limit", 0, "max number of entries to process (0 means all)")
	delay := flag.Duration("delay", 2*time.Second, "delay between fetches")
	userAgent := flag.String("user-agent", defaultUserAgent, "override User-Agent header for fetches")
	cookie := flag.String("cookie", "", "optional Cookie header to include in fetches")
	timeout := flag.Duration("timeout", defaultFetchTimeout, "per-request timeout")
	useMobileHost := flag.Bool("use-mobile-host", true, "rewrite host to m1.codeforces.com for simpler HTML (disable if using browser cookies)")
	flag.Parse()

	problems, err := loadMissingProblems(*csvPath, *start, *limit)
	if err != nil {
		log.Fatalf("failed to load %s: %v", *csvPath, err)
	}
	if len(problems) == 0 {
		log.Println("no problems to process")
		return
	}

	for i, problem := range problems {
		dir := contestPath(problem.ContestID)
		targetDir := filepath.Join(".", dir)
		targetFile := filepath.Join(targetDir, fmt.Sprintf("problem%s.txt", problem.Index))

		log.Printf("[%d/%d] contest %d index %s -> %s", i+1, len(problems), problem.ContestID, problem.Index, targetFile)

		if *dryRun {
			if _, err := os.Stat(targetFile); err == nil {
				log.Printf("  would skip (file exists)")
			} else {
				log.Printf("  would fetch %s", problem.URL)
			}
			continue
		}

		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			log.Printf("  failed to create directory %s: %v", targetDir, err)
			continue
		}

		if _, err := os.Stat(targetFile); err == nil {
			log.Printf("  skipping; %s already exists", targetFile)
			continue
		}

		statement, err := fetchProblemStatement(httpClient, problem.URL, *userAgent, *cookie, *timeout, *useMobileHost)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Printf("  timeout fetching %s: %v", problem.URL, err)
			} else {
				log.Printf("  failed to fetch %s: %v", problem.URL, err)
			}
			continue
		}

		if err := os.WriteFile(targetFile, []byte(statement), 0o644); err != nil {
			log.Printf("  failed to write %s: %v", targetFile, err)
			continue
		}

		log.Printf("  saved %s", targetFile)
		time.Sleep(*delay)
	}
}

func loadMissingProblems(csvPath string, start, limit int) ([]missingProblem, error) {
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}

	var problems []missingProblem
	seen := make(map[string]struct{})
	index := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		index++
		if index <= start {
			continue
		}
		if limit > 0 && len(problems) >= limit {
			break
		}

		contestID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("skipping row with invalid contest id %q: %v", record[0], err)
			continue
		}

		key := fmt.Sprintf("%d-%s", contestID, record[1])
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		problems = append(problems, missingProblem{
			ContestID: contestID,
			Index:     record[1],
			Name:      record[2],
			URL:       record[3],
		})
	}

	return problems, nil
}

func contestPath(contestID int) string {
	thousands := (contestID / 1000) * 1000
	tDir := fmt.Sprintf("%d-%d", thousands, thousands+999)

	remainder := contestID % 1000
	hundreds := (remainder / 100) * 100
	hDirStart := thousands + hundreds
	hDir := fmt.Sprintf("%d-%d", hDirStart, hDirStart+99)

	remainder = remainder % 100
	tens := (remainder / 10) * 10
	teDirStart := thousands + hundreds + tens
	teDir := fmt.Sprintf("%d-%d", teDirStart, teDirStart+9)

	return filepath.Join(tDir, hDir, teDir, fmt.Sprintf("%d", contestID))
}

func fetchProblemStatement(client *http.Client, rawURL, userAgent, cookie string, timeout time.Duration, useMobileHost bool) (string, error) {
	normalizedURL, err := normalizeProblemURL(rawURL, useMobileHost)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", normalizedURL, nil)
	if err != nil {
		return "", err
	}
	if userAgent == "" {
		userAgent = defaultUserAgent
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://codeforces.com/problemset")
	if strings.TrimSpace(cookie) != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", fmt.Errorf("request timed out after %s: %w", timeout, err)
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	selection := doc.Find(".problem-statement")
	if selection.Length() == 0 {
		return "", fmt.Errorf("problem statement not found in response")
	}

	text := strings.TrimSpace(selection.Text())
	if text == "" {
		return "", fmt.Errorf("empty problem statement after parsing")
	}

	return normalizeLineEndings(text), nil
}

func normalizeLineEndings(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func normalizeProblemURL(raw string, useMobileHost bool) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" {
		parsed.Scheme = "https"
	}
	if useMobileHost {
		parsed.Host = "m1.codeforces.com"
	}

	query := parsed.Query()
	if _, ok := query["locale"]; !ok {
		query.Set("locale", "en")
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}
