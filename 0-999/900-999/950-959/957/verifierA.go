package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TestCase struct {
	n        int
	s        string
	expected string
}

func solveCase(n int, s string) string {
	colors := []byte{'C', 'M', 'Y'}
	dp := [3]int{}
	// first character
	for i, col := range colors {
		if s[0] == '?' || s[0] == col {
			dp[i] = 1
		}
	}
	for i := 1; i < n; i++ {
		newDp := [3]int{}
		for j, col := range colors {
			if s[i] != '?' && s[i] != col {
				continue
			}
			for k := 0; k < 3; k++ {
				if k == j {
					continue
				}
				newDp[j] += dp[k]
				if newDp[j] > 2 {
					newDp[j] = 2
				}
			}
		}
		dp = newDp
	}
	total := 0
	for i := 0; i < 3; i++ {
		total += dp[i]
		if total > 2 {
			total = 2
			break
		}
	}
	if total >= 2 {
		return "Yes"
	}
	return "No"
}

func run(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildRandomCase(rng *rand.Rand) TestCase {
	n := rng.Intn(100) + 1
	letters := []byte{'C', 'M', 'Y', '?'}
	sb := strings.Builder{}
	arr := make([]byte, n)
	for i := 0; i < n; i++ {
		arr[i] = letters[rng.Intn(len(letters))]
	}
	sb.WriteString(fmt.Sprintf("%d\n%s\n", n, string(arr)))
	return TestCase{n: n, s: string(arr), expected: solveCase(n, string(arr))}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []TestCase
	// deterministic edge cases
	cases = append(cases, TestCase{n: 1, s: "?", expected: solveCase(1, "?")})
	cases = append(cases, TestCase{n: 1, s: "C", expected: solveCase(1, "C")})
	cases = append(cases, TestCase{n: 2, s: "??", expected: solveCase(2, "??")})
	cases = append(cases, TestCase{n: 3, s: "CM?", expected: solveCase(3, "CM?")})
	cases = append(cases, TestCase{n: 4, s: "C?M?", expected: solveCase(4, "C?M?")})

	for i := 0; i < 100; i++ {
		cases = append(cases, buildRandomCase(rng))
	}

	for i, tc := range cases {
		input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
		expected := tc.expected
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
