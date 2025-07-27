package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	A string
	B string
}

func expectedScore(a, b string) int {
	n, m := len(a), len(b)
	dpPrev := make([]int, m+1)
	dpCur := make([]int, m+1)
	maxScore := 0
	for i := 1; i <= n; i++ {
		dpCur[0] = 0
		for j := 1; j <= m; j++ {
			best := dpPrev[j] - 1
			if v := dpCur[j-1] - 1; v > best {
				best = v
			}
			if a[i-1] == b[j-1] {
				if v := dpPrev[j-1] + 2; v > best {
					best = v
				}
			}
			if best < 0 {
				best = 0
			}
			dpCur[j] = best
			if best > maxScore {
				maxScore = best
			}
		}
		dpPrev, dpCur = dpCur, dpPrev
	}
	return maxScore
}

func runBin(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCase {
	tests := []testCase{
		{A: "abc", B: "abc"},
		{A: "abacaba", B: "acababa"},
		{A: "aaa", B: "bbb"},
		{A: "abcd", B: "bc"},
		{A: "a", B: "a"},
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := r.Intn(10) + 1
		m := r.Intn(10) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(byte('a' + r.Intn(4)))
		}
		A := sb.String()
		sb.Reset()
		for j := 0; j < m; j++ {
			sb.WriteByte(byte('a' + r.Intn(4)))
		}
		B := sb.String()
		tests = append(tests, testCase{A: A, B: B})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n%s\n%s\n", len(t.A), len(t.B), t.A, t.B)
		want := expectedScore(t.A, t.B)
		out, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: cannot parse output\n", i+1)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
