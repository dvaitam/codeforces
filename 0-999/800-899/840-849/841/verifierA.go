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
	n      int
	k      int
	s      string
	expect string
}

func solve(n, k int, s string) string {
	counts := make(map[rune]int)
	for _, ch := range s {
		counts[ch]++
	}
	for _, v := range counts {
		if v > k {
			return "NO"
		}
	}
	return "YES"
}

func generateTests() []testCase {
	r := rand.New(rand.NewSource(42))
	tests := make([]testCase, 0, 100)

	// some fixed edge cases
	tests = append(tests, testCase{1, 1, "a", solve(1, 1, "a")})
	tests = append(tests, testCase{2, 1, "aa", solve(2, 1, "aa")})
	tests = append(tests, testCase{2, 2, "aa", solve(2, 2, "aa")})
	tests = append(tests, testCase{3, 1, "abc", solve(3, 1, "abc")})

	for len(tests) < 100 {
		n := r.Intn(100) + 1
		k := r.Intn(100) + 1
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = byte('a' + r.Intn(26))
		}
		s := string(b)
		tests = append(tests, testCase{n, k, s, solve(n, k, s)})
	}
	return tests
}

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = os.Environ()
	// set timeout
	timer := time.AfterFunc(2*time.Second, func() {
		cmd.Process.Kill()
	})
	err := cmd.Run()
	timer.Stop()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	path := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n%s\n", tc.n, tc.k, tc.s)
		out, err := runBinary(path, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nOutput: %s\n", i+1, err, out)
			os.Exit(1)
		}
		res := strings.TrimSpace(out)
		if res != tc.expect {
			fmt.Printf("Test %d failed\nInput:\n%d %d\n%s\nExpected: %s\nGot: %s\n", i+1, tc.n, tc.k, tc.s, tc.expect, res)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
