package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	missing := make(map[string][]string)

	re := regexp.MustCompile(`sol([A-Z])\.cpp`)
	goRe := regexp.MustCompile(`(\d+)([A-Z])\.go`)

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		base := filepath.Base(dir)

		var letter string
		if match := re.FindStringSubmatch(d.Name()); len(match) > 1 {
			letter = match[1]
		} else if goMatch := goRe.FindStringSubmatch(d.Name()); len(goMatch) > 2 {
			problemNum := goMatch[1]
			if problemNum == base {
				letter = goMatch[2]
			}
		}

		if letter != "" {
			problemStatementPath := filepath.Join(dir, "problem"+letter+".txt")
			if _, err := os.Stat(problemStatementPath); os.IsNotExist(err) {

				if _, ok := missing[base]; !ok {
					missing[base] = []string{}
				}
				// avoid duplicates
				found := false
				for _, l := range missing[base] {
					if l == letter {
						found = true
						break
					}
				}
				if !found {
					missing[base] = append(missing[base], letter)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "error walking directory: %v\n", err)
		os.Exit(1)
	}

	// Create a new browser context
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/chromium-browser"),
		chromedp.Flag("headless", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	for problem, letters := range missing {
		for _, letter := range letters {
			fmt.Printf("Fetching problem statement for problem %s, letter %s\n", problem, letter)
			url := fmt.Sprintf("https://codeforces.com/problemset/problem/%s/%s", problem, letter)

			var problemStatement string
			err := chromedp.Run(taskCtx,
				chromedp.Navigate(url),
				chromedp.WaitVisible(".problem-statement", chromedp.ByQuery),
				chromedp.Text(".problem-statement", &problemStatement, chromedp.ByQuery),
			)

			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to fetch problem statement for %s%s: %v\n", problem, letter, err)
				continue
			}

			// Find the directory for the problem
			dir, err := findProblemDir(problem)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not find directory for problem %s: %v\n", problem, err)
				continue
			}

			// Save the problem statement
			problemStatementPath := filepath.Join(dir, "problem"+letter+".txt")
			err = os.WriteFile(problemStatementPath, []byte(problemStatement), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to write problem statement for %s%s: %v\n", problem, letter, err)
			}

			fmt.Printf("Successfully fetched and saved problem statement for %s%s\n", problem, letter)
			time.Sleep(2 * time.Second) // Be nice to Codeforces
		}
	}
}

func findProblemDir(problemID string) (string, error) {
	var dir string
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == problemID {
			dir = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", fmt.Errorf("directory for problem %s not found", problemID)
	}
	return dir, nil
}