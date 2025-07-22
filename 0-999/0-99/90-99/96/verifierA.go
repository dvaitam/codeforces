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

type testCase struct {
	s      string
	expect string
}

func isDangerous(s string) string {
	if len(s) == 0 {
		return "NO"
	}
	count := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			count++
			if count >= 7 {
				return "YES"
			}
		} else {
			count = 1
		}
	}
	if count >= 7 {
		return "YES"
	}
	return "NO"
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	// deterministic edge cases
	edges := []string{
		"0000000",                      // exactly seven zeros -> YES
		"1111111",                      // seven ones -> YES
		"0101010",                      // alternating -> NO
		"0000001",                      // six zeros -> NO
		"1111110",                      // six ones -> NO
		"0000000000000",                // long zeros -> YES
		strings.Repeat("01", 50)[:100], // max length alternating -> NO
	}
	for _, e := range edges {
		tests = append(tests, testCase{s: e, expect: isDangerous(e)})
	}
	// random tests until we have at least 100
	for len(tests) < 100 {
		l := rng.Intn(100) + 1
		var sb strings.Builder
		for i := 0; i < l; i++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		s := sb.String()
		tests = append(tests, testCase{s: s, expect: isDangerous(s)})
	}
	return tests
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.s + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	result := strings.TrimSpace(out.String())
	if result != tc.expect {
		return fmt.Errorf("expected %s got %s", tc.expect, result)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %s\n", i+1, err, tc.s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
