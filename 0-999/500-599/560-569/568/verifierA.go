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

const N = 2000000

var piCount []int
var rubCount []int

func isPalindrome(n int) bool {
	rev := 0
	tmp := n
	for tmp > 0 {
		rev = rev*10 + tmp%10
		tmp /= 10
	}
	return rev == n
}

func initCounts() {
	primes := make([]bool, N+1)
	for i := 2; i <= N; i++ {
		primes[i] = true
	}
	for i := 2; i*i <= N; i++ {
		if primes[i] {
			for j := i * i; j <= N; j += i {
				primes[j] = false
			}
		}
	}
	piCount = make([]int, N+1)
	rubCount = make([]int, N+1)
	for i := 1; i <= N; i++ {
		piCount[i] = piCount[i-1]
		rubCount[i] = rubCount[i-1]
		if primes[i] {
			piCount[i]++
		}
		if isPalindrome(i) {
			rubCount[i]++
		}
	}
}

func solveCase(p, q int) string {
	ans := -1
	for n := 1; n <= N; n++ {
		if piCount[n]*q <= rubCount[n]*p {
			ans = n
		}
	}
	if ans == -1 {
		return "Palindromic tree is better than splay tree\n"
	}
	return fmt.Sprintf("%d\n", ans)
}

type testCase struct {
	p, q     int
	expected string
}

func generateRandomCase(rng *rand.Rand) testCase {
	q := rng.Intn(10000) + 1
	maxP := q * 42
	if maxP > 10000 {
		maxP = 10000
	}
	minP := (q + 41) / 42
	if minP < 1 {
		minP = 1
	}
	if minP > maxP {
		minP = maxP
	}
	p := rng.Intn(maxP-minP+1) + minP
	return testCase{p: p, q: q, expected: solveCase(p, q)}
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.p, tc.q)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(tc.expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(tc.expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	initCounts()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	// deterministic edge cases
	edges := []struct{ p, q int }{{1, 42}, {42, 1}, {1, 1}, {2, 3}, {10000, 10000}}
	for _, e := range edges {
		cases = append(cases, testCase{p: e.p, q: e.q, expected: solveCase(e.p, e.q)})
	}

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d %d\n", i+1, err, tc.p, tc.q)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
