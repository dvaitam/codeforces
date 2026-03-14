package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// expected computes the correct answer for problem 1990A.
// Alice wins iff there exists some value (1..n) with an odd frequency.
func expected(arr []int) string {
	n := len(arr)
	cnt := make([]int, n+1)
	for _, v := range arr {
		if v >= 1 && v <= n {
			cnt[v]++
		}
	}
	for i := 1; i <= n; i++ {
		if cnt[i]%2 == 1 {
			return "YES"
		}
	}
	return "NO"
}

type testCase struct {
	arr []int
}

// generateTests creates random test cases respecting constraints:
// 2 <= n <= 50, 1 <= a_i <= n.
func generateTests(numTests int) []testCase {
	tests := make([]testCase, numTests)
	for i := 0; i < numTests; i++ {
		n := rand.Intn(49) + 2 // n in [2, 50]
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(n) + 1 // a_j in [1, n]
		}
		tests[i] = testCase{arr: arr}
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.arr)))
		for _, v := range tc.arr {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	tests := generateTests(100)
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	lines := strings.Fields(output)
	if len(lines) != len(tests) {
		fmt.Printf("expected %d outputs, got %d\n", len(tests), len(lines))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := expected(tc.arr)
		got := strings.ToUpper(strings.TrimSpace(lines[i]))
		if got != want {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, want, lines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
