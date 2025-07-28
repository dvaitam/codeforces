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
	input    string
	expected string
}

func cmp(a, b byte) int8 {
	if a < b {
		return 1
	} else if a > b {
		return -1
	}
	return 0
}

func solveGame(s string) string {
	n := len(s)
	dp := make([][]int8, n)
	for i := range dp {
		dp[i] = make([]int8, n)
	}
	for i := 0; i+1 < n; i++ {
		if s[i] == s[i+1] {
			dp[i][i+1] = 0
		} else {
			dp[i][i+1] = 1
		}
	}
	for length := 4; length <= n; length += 2 {
		for l := 0; l+length-1 < n; l++ {
			r := l + length - 1
			a1 := dp[l+2][r]
			if a1 == 0 {
				a1 = cmp(s[l], s[l+1])
			}
			a2 := dp[l+1][r-1]
			if a2 == 0 {
				a2 = cmp(s[l], s[r])
			}
			if a1 > a2 {
				a1 = a2
			}
			b1 := dp[l][r-2]
			if b1 == 0 {
				b1 = cmp(s[r], s[r-1])
			}
			b2 := dp[l+1][r-1]
			if b2 == 0 {
				b2 = cmp(s[r], s[l])
			}
			if b1 > b2 {
				b1 = b2
			}
			if a1 < b1 {
				dp[l][r] = b1
			} else {
				dp[l][r] = a1
			}
		}
	}
	res := dp[0][n-1]
	if res > 0 {
		return "Alice"
	} else if res < 0 {
		return "Bob"
	}
	return "Draw"
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(4))
	tests := []testCase{}
	tests = append(tests, testCase{input: "1\naa\n", expected: "Draw"})
	tests = append(tests, testCase{input: "1\naba\n", expected: "Alice"})
	for len(tests) < 100 {
		length := (rng.Intn(3) + 1) * 2
		b := make([]byte, length)
		for i := 0; i < length; i++ {
			b[i] = byte('a' + rng.Intn(3))
		}
		s := string(b)
		expect := solveGame(s)
		tests = append(tests, testCase{input: fmt.Sprintf("1\n%s\n", s), expected: expect})
	}
	return tests
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_ = rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := generateTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
