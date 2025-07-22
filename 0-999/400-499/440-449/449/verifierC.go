package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type pair struct{ a, b int }

type test struct {
	input, expected string
}

func solve(input string) string {
	var n int
	fmt.Sscan(strings.TrimSpace(input), &n)
	used := make([]bool, n+1)
	grp := make([]pair, 0, n/2)
	var p int
	for i := 3; i <= n/2; i += 2 {
		if !used[i] {
			p = i
			for j := i * 3; j <= n; j += i {
				if !used[j] {
					if p != 0 {
						grp = append(grp, pair{p, j})
						used[p], used[j] = true, true
						p = 0
					} else {
						p = j
					}
				}
			}
			if p != 0 {
				pairv := i * 2
				if pairv <= n && !used[p] && !used[pairv] {
					grp = append(grp, pair{p, pairv})
					used[p], used[pairv] = true, true
				}
			}
		}
	}
	p = 0
	for i := 2; i <= n; i += 2 {
		if !used[i] {
			if p != 0 {
				grp = append(grp, pair{p, i})
				used[p], used[i] = true, true
				p = 0
			} else {
				p = i
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(grp))
	for _, pr := range grp {
		fmt.Fprintf(&sb, "%d %d\n", pr.a, pr.b)
	}
	return strings.TrimSpace(sb.String())
}

func generateTests() []test {
	rand.Seed(451)
	var tests []test
	fixed := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 50}
	for _, n := range fixed {
		inp := fmt.Sprintf("%d\n", n)
		tests = append(tests, test{inp, solve(inp)})
	}
	for len(tests) < 100 {
		n := rand.Intn(100) + 1
		inp := fmt.Sprintf("%d\n", n)
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != t.expected {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
