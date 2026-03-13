package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// 825D: Suitable Replacement
// Replace '?' in s with lowercase letters to maximize suitability (number of
// non-overlapping occurrences of t after arbitrary rearrangement of s).
// Multiple valid outputs exist; we validate by checking suitability matches optimal.

func computeSuitability(s, t string) int {
	// After rearranging s, max non-overlapping copies of t.
	// Count chars in s and t, then greedily fit copies.
	if len(t) == 0 {
		return 0
	}
	sc := make([]int, 26)
	for i := 0; i < len(s); i++ {
		sc[int(s[i]-'a')]++
	}
	tc := make([]int, 26)
	for i := 0; i < len(t); i++ {
		tc[int(t[i]-'a')]++
	}
	// Binary search on number of copies
	lo, hi := 0, len(s)
	for lo < hi {
		mid := (lo + hi + 1) / 2
		ok := true
		for c := 0; c < 26; c++ {
			if mid*tc[c] > sc[c] {
				ok = false
				break
			}
		}
		if ok {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func optimalSuitability(s, t string) int {
	// Compute maximum possible suitability over all replacements of '?' in s
	if len(t) == 0 {
		return 0
	}
	sc := make([]int, 26)
	qCount := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '?' {
			qCount++
		} else {
			sc[int(s[i]-'a')]++
		}
	}
	tc := make([]int, 26)
	for i := 0; i < len(t); i++ {
		tc[int(t[i]-'a')]++
	}

	lo, hi := 0, len(s)
	for lo < hi {
		mid := (lo + hi + 1) / 2
		need := 0
		for c := 0; c < 26; c++ {
			deficit := mid*tc[c] - sc[c]
			if deficit > 0 {
				need += deficit
			}
		}
		if need <= qCount {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func randStringLetters(rng *rand.Rand, n int, allowQ bool) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if allowQ && rng.Intn(5) == 0 {
			b[i] = '?'
		} else {
			b[i] = byte('a' + rng.Intn(26))
		}
	}
	return string(b)
}

func genCase(rng *rand.Rand) (string, string, string) {
	n := rng.Intn(30) + 1
	m := rng.Intn(15) + 1
	s := randStringLetters(rng, n, true)
	t := randStringLetters(rng, m, false)
	input := fmt.Sprintf("%s\n%s\n", s, t)
	return input, s, t
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, s, t := genCase(rng)
		optSuit := optimalSuitability(s, t)

		got, err := run(exe, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", input)
			os.Exit(1)
		}

		// Validate candidate output
		gotStr := strings.TrimSpace(got)

		// Check length matches
		if len(gotStr) != len(s) {
			fmt.Printf("wrong answer on test %d: output length %d != input length %d\n", i+1, len(gotStr), len(s))
			fmt.Println("input:\n", input)
			fmt.Println("got:\n", gotStr)
			os.Exit(1)
		}

		// Check no '?' in output
		if strings.Contains(gotStr, "?") {
			fmt.Printf("wrong answer on test %d: output contains '?'\n", i+1)
			fmt.Println("input:\n", input)
			fmt.Println("got:\n", gotStr)
			os.Exit(1)
		}

		// Check all chars are lowercase letters
		allLower := true
		for _, c := range gotStr {
			if c < 'a' || c > 'z' {
				allLower = false
				break
			}
		}
		if !allLower {
			fmt.Printf("wrong answer on test %d: output contains non-lowercase chars\n", i+1)
			fmt.Println("input:\n", input)
			fmt.Println("got:\n", gotStr)
			os.Exit(1)
		}

		// Check that non-'?' positions in s are preserved
		for j := 0; j < len(s); j++ {
			if s[j] != '?' && gotStr[j] != s[j] {
				fmt.Printf("wrong answer on test %d: position %d changed from '%c' to '%c'\n", i+1, j, s[j], gotStr[j])
				fmt.Println("input:\n", input)
				fmt.Println("got:\n", gotStr)
				os.Exit(1)
			}
		}

		// Check suitability matches optimal
		gotSuit := computeSuitability(gotStr, t)
		if gotSuit != optSuit {
			fmt.Printf("wrong answer on test %d: suitability %d, optimal %d\n", i+1, gotSuit, optSuit)
			fmt.Println("input:\n", input)
			fmt.Println("got:\n", gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
