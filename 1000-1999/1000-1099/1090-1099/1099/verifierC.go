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

type unit struct {
	char rune
	mod  rune // 0, '?', or '*'
}

func parseUnits(s string) []unit {
	var parsed []unit
	raw := []rune(s)
	n := len(raw)
	for i := 0; i < n; i++ {
		u := unit{char: raw[i], mod: 0}
		if i+1 < n {
			if raw[i+1] == '?' || raw[i+1] == '*' {
				u.mod = raw[i+1]
				i++
			}
		}
		parsed = append(parsed, u)
	}
	return parsed
}

func getBounds(units []unit) (int, int) {
	minLen := 0
	hasSnow := false
	for _, u := range units {
		if u.mod == 0 {
			minLen++
		}
		if u.mod == '*' {
			hasSnow = true
		}
	}
	maxLen := len(units) // count of base chars
	if hasSnow {
		maxLen = 1000000 // effectively infinite
	}
	return minLen, maxLen
}

func isValid(s string, k int, out string) bool {
	units := parseUnits(s)
	minLen, maxLen := getBounds(units)

	// Logic check for "Impossible"
	isPossible := k >= minLen && k <= maxLen
	
	if !isPossible {
		return out == "Impossible"
	}

	if out == "Impossible" {
		// It was possible but candidate said Impossible
		return false
	}

	if len(out) != k {
		return false
	}

	// DP check
	m := len(units)
	n := len(out)
	outRunes := []rune(out)

	// dp[i][j] = can units[i:] match out[j:]
	dp := make([][]bool, m+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	dp[m][n] = true

	for i := m - 1; i >= 0; i-- {
		u := units[i]
		for j := n; j >= 0; j-- {
			// Option 1: Skip unit (for ? and *)
			if (u.mod == '?' || u.mod == '*') && dp[i+1][j] {
				dp[i][j] = true
				continue
			}

			// Option 2: Match character
			if j < n && outRunes[j] == u.char {
				if u.mod == 0 {
					if dp[i+1][j+1] {
						dp[i][j] = true
					}
				} else if u.mod == '?' {
					if dp[i+1][j+1] {
						dp[i][j] = true
					}
				} else if u.mod == '*' {
					// For *, we can match char and stay at i (to match more)
					// or move to i+1 (handled by skip logic via recursion if we were doing recursion,
					// but here we need to explicitely link to states)
					
					// Actually, dp[i][j] = skip || (match && dp[i][j+1])
					// "skip" corresponds to matching 0.
					// "match && dp[i][j+1]" corresponds to matching 1+, because from dp[i][j+1] we can eventually skip.
					if dp[i][j+1] {
						dp[i][j] = true
					}
				}
			}
		}
	}

	return dp[0][0]
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 100)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 100; i++ {
		l := rng.Intn(15) + 1
		var sb strings.Builder
		for j := 0; j < l; j++ {
			ch := letters[rng.Intn(len(letters))]
			sb.WriteRune(ch)
			r := rng.Intn(3)
			if r == 0 {
				sb.WriteByte('?')
			} else if r == 1 {
				sb.WriteByte('*')
			}
		}
		k := rng.Intn(20) + 1
		tests[i] = fmt.Sprintf("%s\n%d\n", sb.String(), k)
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	tests := genTests()
	// Add the failing test case manually to ensure it passes
	tests = append(tests, "q?u*s*r?e?d*w?t?o*\n20\n")

	for idx, input := range tests {
		out, err := runBinary(cand, input)
		if err != nil {
			fmt.Printf("runtime error on test %d\n", idx+1)
			os.Exit(1)
		}
		
		lines := strings.Split(strings.TrimSpace(input), "\n")
		s := lines[0]
		var k int
		fmt.Sscanf(lines[1], "%d", &k)

		if !isValid(s, k, out) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sgot:\n%s\n", idx+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}