package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

func runCandidate(binary string, input []byte) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\n"), nil
}

// parseNames extracts the list of names from test input.
func parseNames(input []byte) []string {
	lines := strings.Split(strings.TrimRight(string(input), "\n"), "\n")
	if len(lines) == 0 {
		return nil
	}
	var n int
	fmt.Sscanf(lines[0], "%d", &n)
	names := make([]string, 0, n)
	for i := 1; i <= n && i < len(lines); i++ {
		names = append(names, strings.TrimSpace(lines[i]))
	}
	return names
}

// isSortedUnder checks that names are lexicographically sorted under the given
// alphabet order (a 26-char permutation of 'a'..'z').
// Returns an error describing the first violation, or nil if valid.
func isSortedUnder(names []string, order string) error {
	rank := [26]int{}
	for i, c := range order {
		rank[c-'a'] = i
	}
	for i := 0; i+1 < len(names); i++ {
		s1, s2 := names[i], names[i+1]
		minLen := len(s1)
		if len(s2) < minLen {
			minLen = len(s2)
		}
		diff := false
		for j := 0; j < minLen; j++ {
			r1, r2 := rank[s1[j]-'a'], rank[s2[j]-'a']
			if r1 < r2 {
				diff = true
				break
			}
			if r1 > r2 {
				return fmt.Errorf("names[%d]=%q comes after names[%d]=%q under given order (char %q > %q)", i, s1, i+1, s2, s1[j], s2[j])
			}
		}
		if !diff && len(s1) > len(s2) {
			return fmt.Errorf("names[%d]=%q is longer prefix of names[%d]=%q but listed first", i, s1, i+1, s2)
		}
	}
	return nil
}

// isValidPermutation returns true if s is exactly the 26 lowercase letters each once.
func isValidPermutation(s string) bool {
	if len(s) != 26 {
		return false
	}
	seen := [26]bool{}
	for _, c := range s {
		if c < 'a' || c > 'z' {
			return false
		}
		if seen[c-'a'] {
			return false
		}
		seen[c-'a'] = true
	}
	return true
}

func runTests(dir, binary string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.in"))
	if err != nil {
		return err
	}
	sort.Strings(files)
	for _, inFile := range files {
		outFile := inFile[:len(inFile)-3] + ".out"
		input, err := os.ReadFile(inFile)
		if err != nil {
			return err
		}
		expected, err := os.ReadFile(outFile)
		if err != nil {
			return err
		}
		expectedStr := strings.TrimRight(string(expected), "\n")

		got, err := runCandidate(binary, input)
		if err != nil {
			return fmt.Errorf("%s: %v", filepath.Base(inFile), err)
		}

		base := filepath.Base(inFile)

		if expectedStr == "Impossible" {
			// Exact match required
			if got != "Impossible" {
				return fmt.Errorf("%s: expected Impossible but got %q", base, got)
			}
			continue
		}

		// Expected is a valid permutation — accept any correct permutation.
		if got == "Impossible" {
			return fmt.Errorf("%s: got Impossible but a valid ordering exists", base)
		}
		if !isValidPermutation(got) {
			return fmt.Errorf("%s: output %q is not a valid permutation of a-z", base, got)
		}
		names := parseNames(input)
		if err := isSortedUnder(names, got); err != nil {
			return fmt.Errorf("%s: %v", base, err)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Dir(file)
	testDir := filepath.Join(base, "tests", "C")
	if err := runTests(testDir, binary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
