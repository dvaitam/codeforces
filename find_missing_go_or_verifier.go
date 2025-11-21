package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

type problemReport struct {
	ContestID   string
	ContestPath string
	Index       string
	HasSolution bool
	HasVerifier bool
}

func main() {
	reports, err := collectReports(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning contests: %v\n", err)
		os.Exit(1)
	}

	if len(reports) == 0 {
		fmt.Println("Every problem with a statement already has both a Go solution and a verifier.")
		return
	}

	sort.Slice(reports, func(i, j int) bool {
		if reports[i].ContestID != reports[j].ContestID {
			return reports[i].ContestID < reports[j].ContestID
		}
		if reports[i].ContestPath != reports[j].ContestPath {
			return reports[i].ContestPath < reports[j].ContestPath
		}
		return reports[i].Index < reports[j].Index
	})

	for _, r := range reports {
		var missing []string
		if !r.HasSolution {
			missing = append(missing, "Go solution")
		}
		if !r.HasVerifier {
			missing = append(missing, "verifier")
		}
		fmt.Printf("%s%s in %s is missing %s\n",
			r.ContestID, r.Index, r.ContestPath, strings.Join(missing, " and "))
	}
}

func collectReports(root string) ([]problemReport, error) {
	var reports []problemReport

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if !isContestDir(d.Name()) {
			return nil
		}

		contestReports, err := inspectContest(path, d.Name())
		if err != nil {
			return err
		}
		reports = append(reports, contestReports...)
		return filepath.SkipDir
	})

	return reports, err
}

func inspectContest(dir, contestID string) ([]problemReport, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var reports []problemReport
	indexes := findProblemIndexes(entries)
	if len(indexes) == 0 {
		return nil, nil
	}

	generalVerifier := fileExists(filepath.Join(dir, "verifier.go"))

	for _, idx := range indexes {
		solution := fmt.Sprintf("%s%s.go", contestID, idx)
		hasSolution := fileExists(filepath.Join(dir, solution))

		hasVerifier := generalVerifier
		if !hasVerifier {
			verifier := fmt.Sprintf("verifier%s.go", idx)
			hasVerifier = fileExists(filepath.Join(dir, verifier))
		}

		if hasSolution && hasVerifier {
			continue
		}

		reports = append(reports, problemReport{
			ContestID:   contestID,
			ContestPath: dir,
			Index:       idx,
			HasSolution: hasSolution,
			HasVerifier: hasVerifier,
		})
	}

	return reports, nil
}

func findProblemIndexes(entries []os.DirEntry) []string {
	seen := make(map[string]struct{})
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, "problem") || !strings.HasSuffix(name, ".txt") {
			continue
		}

		idx := strings.TrimSuffix(strings.TrimPrefix(name, "problem"), ".txt")
		if idx == "" {
			continue
		}

		seen[idx] = struct{}{}
	}

	indexes := make([]string, 0, len(seen))
	for idx := range seen {
		indexes = append(indexes, idx)
	}

	sort.Strings(indexes)
	return indexes
}

func isContestDir(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
